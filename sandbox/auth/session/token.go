package session

import (
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"golang.org/x/oauth2"
)

// Tokens can generally be ID or access. It's a bit weird to have a struct
// named "Token" that only contains an ID token. An alternative is to just rename
// this struct to IDToken.
type Token struct {
	// Are these fields needed?
	Issuer string            `json:"issuer,omitempty"`
	Keys   []jose.JSONWebKey `json:"keys,omitempty"`

	IDToken string `json:"id_token,omitempty"`

	idClaims *IDClaims
}

func TokenFromString(t string) (*Token, error) {
	_, err := jwt.ParseSigned(t)
	if err != nil {
		return nil, err
	}
	tok := &Token{
		IDToken: t,
	}
	claims := tok.IDClaims()
	if claims != nil {
		tok.Issuer = claims.Issuer
	}
	return tok, nil
}

func (t *Token) ToOauth2Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: t.IDToken,
	}
}
