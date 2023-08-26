package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang-jwt/jwt/v5"
	"go.jetpack.io/envsec/debug"
)

// for valiation of JWT (can be hardcoded to avoid fetching)
const jwksURL = `https://jetpack-io.us.auth0.com/.well-known/jwks.json`
const awsRegion = "us-east-1"

// have to be var instead of const so that their &address can be passed
var accountId = "984256416385"
var identityPoolId = "us-east-1:da3c3c71-61c7-4f7c-8e3d-3770e9b61379"

// Data that will be in JWT
type UserClaim struct {
	jwt.RegisteredClaims
	OrgID string `json:"https://auth.jetpack.io/org_id,omitempty"`
}

// getting the saved auth token after login
func getAuthToken() (string, error) {
	type auth struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IdToken      string `json:"id_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	authFile, err := os.Open(fmt.Sprintf("%s/.local/state/devbox/auth.json", homedir))
	if err != nil {
		return "", err
	}
	defer authFile.Close()

	byteContent, err := io.ReadAll(authFile)
	if err != nil {
		return "", err
	}

	var result auth
	json.Unmarshal([]byte(byteContent), &result)
	return result.AccessToken, nil
}

// fetching orgID from accessToken
func getOrgId() (string, error) {
	authToken, err := getAuthToken()
	if err != nil {
		return "", err
	}
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return "", err
	}
	// Parse the token with custom orgId claim
	var userClaim UserClaim
	token, err := jwt.ParseWithClaims(authToken, &userClaim, jwks.Keyfunc)
	if err != nil {
		return "", err
	}

	debug.Log("token claims %v \n", token.Claims.(*UserClaim))
	return token.Claims.(*UserClaim).OrgID, nil

}

func GetAWSCreds() (*cognitoidentity.GetCredentialsForIdentityOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		return nil, err
	}
	svc := cognitoidentity.New(sess)

	auth0Token, err := getAuthToken()
	if err != nil {
		return nil, err
	}
	debug.Log("auth0Token: %s \n\n", auth0Token)
	logins := map[string]*string{
		"auth.jetpack.io": &auth0Token,
	}
	getIdoutput, err := svc.GetId(&cognitoidentity.GetIdInput{
		AccountId:      &accountId,
		IdentityPoolId: &identityPoolId,
		Logins:         logins,
	})
	if err != nil {
		return nil, err
	}
	debug.Log("cognito ID: %v\n\n", getIdoutput.GoString())

	output, err := svc.GetCredentialsForIdentity(&cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: getIdoutput.IdentityId,
		Logins:     logins,
	})
	if err != nil {
		return nil, err
	}
	debug.Log("aws credentials: %v \n\n", output.GoString())
	return output, nil
}

func GetParameters() ([]*ssm.Parameter, error) {
	stsoutput, err := GetAWSCreds()
	if err != nil {
		return nil, err
	}
	authsess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			*stsoutput.Credentials.AccessKeyId,
			*stsoutput.Credentials.SecretKey,
			*stsoutput.Credentials.SessionToken,
		),
	})
	if err != nil {
		return nil, err
	}
	orgId, err := getOrgId()
	if err != nil {
		return nil, err
	}
	ssmsvc := ssm.New(authsess)

	params, err := ssmsvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path: aws.String(fmt.Sprintf("/jetpackio/secrets/%s/", orgId)),
	})

	if err != nil {
		return nil, err
	}

	debug.Log("params: %v\n", params.GoString())
	return params.Parameters, nil

}
