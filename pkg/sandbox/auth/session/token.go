package session

import (
	"github.com/go-jose/go-jose/v3"
	"golang.org/x/oauth2"
)

type Token struct {
	oauth2.Token // Embed and oauth2 token

	Issuer string            `json:"issuer,omitempty"`
	Keys   []jose.JSONWebKey `json:"keys,omitempty"`

	IDToken string `json:"id_token,omitempty"`

	idClaims *IDClaims
}
