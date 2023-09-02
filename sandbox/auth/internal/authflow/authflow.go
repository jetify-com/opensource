package authflow

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/browser"
	"go.jetpack.io/auth/internal/pkce"
	"go.jetpack.io/auth/session"
	"golang.org/x/oauth2"
)

type AuthFlow struct {
	URL string

	issuer   string
	clientID string

	pkceCodeVerifier string
	xsrfState        string
	oidcNonce        string

	oauth2Conf   *oauth2.Config
	oidcProvider *oidc.Provider
}

func New(issuer string, clientID string) (*AuthFlow, error) {
	// TODO: We currently default to using the Auth flow with PCKE
	// we could instead check if the issuer supports the device flow, if
	// it does, use that, otherwise use the PKCE flow.
	flow := &AuthFlow{
		issuer:   issuer,
		clientID: clientID,

		pkceCodeVerifier: pkce.GenerateVerifier(),
		xsrfState:        pkce.GenerateVerifier(),
		oidcNonce:        pkce.GenerateVerifier(),
	}

	err := flow.init()
	if err != nil {
		return nil, err
	}

	return flow, nil
}

func (f *AuthFlow) OpenBrowser() error {
	return browser.OpenURL(f.URL)
}

func (f *AuthFlow) Exchange(code string) (*session.Token, error) {
	ctx := context.Background()
	otok, err := f.oauth2Conf.Exchange(ctx, code, pkce.VerifierOption(f.pkceCodeVerifier))
	if err != nil {
		return nil, err
	}

	tok, err := session.FromOauth2(otok)
	if err != nil {
		return nil, err
	}

	// err = f.verify(tok)
	// if err != nil {
	// 	return nil, err
	// }
	return tok, nil
}
