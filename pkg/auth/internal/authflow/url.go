package authflow

import (
	"strings"

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
	if len(f.audience) > 0 {
		// Audience is not a standard OAuth2 parameter, we follow ory/hydra's convention
		// https://www.ory.sh/docs/hydra/guides/audiences#audience-in-authorization-code-implicit-and-hybrid-flows
		opts = append(opts, oauth2.SetAuthURLParam("audience", strings.Join(f.audience, " ")))
	}
	return opts
}
