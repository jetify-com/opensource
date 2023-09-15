package session

import (
	"github.com/go-jose/go-jose/v3"
	"golang.org/x/oauth2"
)

type Token struct {
	oauth2.Token // Embed an oauth2 token

	Issuer string            `json:"issuer,omitempty"`
	Keys   []jose.JSONWebKey `json:"keys,omitempty"`

	// The id token is technically contained in the original oauth2.Token (inside)
	// of extras. However, the extras don't get serialized to JSON, so the id token
	// is lost of we save the token to disk and read it back. We copy the IDToken
	// into this field so that it does get serialized properly.
	IDToken string `json:"id_token,omitempty"`

	idClaims *IDClaims
}
