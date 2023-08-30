package token

import "time"

type Token struct {
	Issuer       string `json:"issuer"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`

	TokenType string    `json:"token_type,omitempty"`
	Expiry    time.Time `json:"expiry,omitempty"`
}
