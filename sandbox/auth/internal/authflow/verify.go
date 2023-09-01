package authflow

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

///////
//////

func (f *AuthFlow) Verify(tok *oauth2.Token) error {
	ctx := context.Background()
	if !tok.Valid() {
		return fmt.Errorf("oauth token is not valid")
	}

	provider := f.getOIDCProvider()

	rawIDToken, ok := tok.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("oauth token does not contain an id_token")
	}

	idConfig := &oidc.Config{
		ClientID: f.clientID,
	}

	idVerifier := provider.Verifier(idConfig)
	idTok, err := idVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return err
	}

	if idTok.Nonce != f.oidcNonce {
		return fmt.Errorf("invalid id token: nonce mismatch")
	}

	return nil
}
