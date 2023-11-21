package authflow

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"go.jetpack.io/pkg/auth/internal/pkce"
	"golang.org/x/oauth2"
)

func (f *AuthFlow) authURL() string {
	params := f.authURLParams()
	conf := f.getOauth2Conf()
	return conf.AuthCodeURL(f.xsrfState, params...)
}

func (f *AuthFlow) authURLParams() []oauth2.AuthCodeOption {
	opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline, oidc.Nonce(f.oidcNonce)}
	opts = append(opts, pkce.S256ChallengeOption(f.pkceCodeVerifier)...)
	return opts
}
