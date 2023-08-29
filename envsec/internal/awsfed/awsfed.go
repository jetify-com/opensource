package awsfed

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"go.jetpack.io/envsec"
	"go.jetpack.io/envsec/internal/debug"
)

type AWSFed struct {
	JwksURL        string
	Region         string
	AccountId      string
	IdentityPoolId string
}

// Data that will be in JWT

func NewAWSFed() *AWSFed {
	return &AWSFed{
		// for valiation of JWT (can be cached to improve performance)
		JwksURL: "https://jetpack-io.us.auth0.com/.well-known/jwks.json",
		// todo change these values below to resources created by terraform
		Region:         "us-east-1",
		AccountId:      "984256416385",
		IdentityPoolId: "us-east-1:da3c3c71-61c7-4f7c-8e3d-3770e9b61379",
	}
}

func (a *AWSFed) getAWSCreds(accessToken string) (*cognitoidentity.Credentials, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.Region)},
	)
	if err != nil {
		return nil, err
	}
	svc := cognitoidentity.New(sess)

	debug.Log("auth0Token: %s \n\n", accessToken)
	logins := map[string]*string{
		"auth.jetpack.io": &accessToken,
	}
	getIdoutput, err := svc.GetId(&cognitoidentity.GetIdInput{
		AccountId:      &a.AccountId,
		IdentityPoolId: &a.IdentityPoolId,
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
	return output.Credentials, nil
}

func (a *AWSFed) GetSSMConfig(accessToken string) (*envsec.SSMConfig, error) {
	stsCredentials, err := a.getAWSCreds(accessToken)
	if err != nil {
		return nil, err
	}

	return &envsec.SSMConfig{
		Region:          a.Region,
		AccessKeyId:     *stsCredentials.AccessKeyId,
		SecretAccessKey: *stsCredentials.SecretKey,
		SessionToken:    *stsCredentials.SessionToken,
	}, nil

}
