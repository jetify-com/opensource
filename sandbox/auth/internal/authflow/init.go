package authflow

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func (f *AuthFlow) init() error {
	err := f.initOIDCProvider()
	if err != nil {
		return err
	}

	err = f.initOauth2Conf()
	if err != nil {
		return err
	}

	// Everything must be initialized before we can generate the URL.
	f.URL = f.authURL()
	return nil
}

func (f *AuthFlow) getOauth2Conf() *oauth2.Config {
	if f.oauth2Conf == nil {
		err := f.initOauth2Conf()
		if err != nil {
			// This should never happen. When the flow first get created, we initialize all the members.
			// If the initialization fails, we return an error at the time of creation.
			return nil
		}
	}
	return f.oauth2Conf
}

func (f *AuthFlow) initOauth2Conf() error {
	if f.oauth2Conf != nil {
		return nil
	}

	provider := f.getOIDCProvider()

	conf := &oauth2.Config{
		ClientID: f.clientID,
		// TODO: should the scopes be configurable?
		Scopes:   []string{"openid", "offline_access"},
		Endpoint: provider.Endpoint(),
	}
	f.oauth2Conf = conf
	return nil
}

func (f *AuthFlow) getOIDCProvider() *oidc.Provider {
	if f.oidcProvider == nil {
		err := f.initOIDCProvider()
		// This should never happen. When the flow first get created, we initialize all the members.
		// If the initialization fails, we return an error at the time of creation.
		if err != nil {
			return nil
		}
	}
	return f.oidcProvider
}

func (f *AuthFlow) initOIDCProvider() error {
	if f.oidcProvider != nil {
		return nil
	}

	ctx := context.Background()
	// Makes a network request to the discovery endpoint.
	// TODO: cache providers and their JWKS keys, so that we don't
	// need to make a network request every time we create a new flow.
	provider, err := oidc.NewProvider(ctx, f.issuer)
	if err != nil {
		return err
	}
	f.oidcProvider = provider
	return nil
}
