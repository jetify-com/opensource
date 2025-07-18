package session

import (
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// Standard claims:
// oidc: https://openid.net/specs/openid-connect-basic-1_0.html#IDToken
// jwt: https://datatracker.ietf.org/doc/html/rfc7519#section-4
// oidc user: https://openid.net/specs/openid-connect-basic-1_0.html#StandardClaims

type IDClaims struct {
	ID string `json:"jti,omitempty"`

	Issuer    string           `json:"iss,omitempty"`
	Subject   string           `json:"sub,omitempty"`
	Audience  jwt.Audience     `json:"aud,omitempty"`
	Expiry    *jwt.NumericDate `json:"exp,omitempty"`
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`
	IssuedAt  *jwt.NumericDate `json:"iat,omitempty"`

	Nonce           string `json:"nonce,omitempty"`
	AccessTokenHash string `json:"at_hash,omitempty"`

	// Email scope:

	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`

	// Profile scope:

	Name       string           `json:"name,omitempty"`
	GivenName  string           `json:"given_name,omitempty"`
	FamilyName string           `json:"family_name,omitempty"`
	UpdatedAt  *jwt.NumericDate `json:"updated_at,omitempty"`

	// Not a default claim in the standards, but used by Auth0
	// and others.
	OrgID string `json:"org_id,omitempty"`
}

func (t *Token) IDClaims() *IDClaims {
	if t.idClaims != nil {
		return t.idClaims
	}

	if t.IDToken == "" {
		return nil
	}

	// Parse the JWT to validate its structure and ensure it uses a supported algorithm.
	// We then use UnsafeClaimsWithoutVerification to extract claims without signature
	// verification. Per JWT spec (RFC 7519), we accept:
	// - MUST: HS256 ("none" is not supported by go-jose for security reasons)
	// - RECOMMENDED: RS256, ES256
	// - Optional but common: other variants for broader compatibility
	jwtTok, err := jwt.ParseSigned(t.IDToken, []jose.SignatureAlgorithm{
		jose.HS256, jose.HS384, jose.HS512,
		jose.RS256, jose.RS384, jose.RS512,
		jose.ES256, jose.ES384, jose.ES512,
		jose.PS256, jose.PS384, jose.PS512,
		jose.EdDSA,
	})
	if err != nil {
		return nil
	}

	claims := IDClaims{}
	err = jwtTok.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return nil
	}

	return &claims
}
