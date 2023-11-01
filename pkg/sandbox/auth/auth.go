package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.jetpack.io/pkg/sandbox/auth/session"
	"golang.org/x/oauth2"

	"go.jetpack.io/pkg/sandbox/auth/internal/authflow"
	"go.jetpack.io/pkg/sandbox/auth/internal/callbackserver"
	"go.jetpack.io/pkg/sandbox/auth/internal/tokenstore"
)

var ErrNotLoggedIn = fmt.Errorf("not logged in")

type Client struct {
	issuer          string
	clientID        string
	store           *tokenstore.Store
	postLoginHook   hook
	postRefreshHook hook
}

type hook func(context.Context, *session.Token) (*session.Token, error)

func NewClient(
	issuer, clientID string,
	postLoginHook, postRefreshHook hook,
) (*Client, error) {
	store, err := tokenstore.New(storeDir())
	if err != nil {
		return nil, err
	}

	return &Client{
		issuer:          issuer,
		clientID:        clientID,
		store:           store,
		postLoginHook:   postLoginHook,
		postRefreshHook: postRefreshHook,
	}, nil
}

func storeDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "~/.cache"
	}
	return filepath.Join(cacheDir, "jetpack.io", "auth")
}

func (c *Client) LoginFlow(ctx context.Context) (*session.Token, error) {
	tok, err := login(c.issuer, c.clientID)
	if err != nil {
		return nil, err
	}

	if c.postLoginHook != nil {
		tok, err = c.postLoginHook(ctx, tok)
		if err != nil {
			return nil, err
		}
	}

	_ = c.store.WriteToken(c.issuer, c.clientID, tok)
	return tok, nil
}

func (c *Client) LogoutFlow() error {
	// For now we just delete the token from the store.
	// But in the future we might want to revoke the token with the server, and do the oauth logout flow.
	return c.RevokeSession()
}

// GetSession returns the current valid session token, if any. If token is
// expired, it will attempt to refresh it. If no token is found, or is unable
// to be refreshed, it will return error.
func (c *Client) GetSession(ctx context.Context) (*session.Token, error) {
	tok, err := c.store.ReadToken(c.issuer, c.clientID)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotLoggedIn
	} else if err != nil {
		return nil, err
	}

	// Refresh if the token is no longer valid:
	if !tok.Valid() {
		tok, err = c.refresh(ctx, tok)
		if err != nil {
			return nil, err
		}
	}

	return tok, nil
}

func (c *Client) refresh(
	ctx context.Context,
	tok *session.Token,
) (*session.Token, error) {
	// TODO: figure out how to share oidc provider and oauth2 client
	// with auth flow:
	provider, err := oidc.NewProvider(ctx, c.issuer)
	if err != nil {
		return tok, err
	}

	conf := oauth2.Config{
		ClientID: c.clientID,
		Endpoint: provider.Endpoint(),
		Scopes:   []string{"openid", "offline_access"},
	}

	// Refresh logic:
	tokenSource := conf.TokenSource(ctx, &tok.Token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return tok, err
	}

	tok.Token = *newToken
	idToken, ok := newToken.Extra("id_token").(string)
	if !ok {
		return tok, fmt.Errorf("missing id_token")
	}
	tok.IDToken = idToken
	if c.postRefreshHook != nil {
		tok, err = c.postLoginHook(ctx, tok)
		if err != nil {
			return nil, err
		}
	}
	err = c.store.WriteToken(c.issuer, c.clientID, tok)
	if err != nil {
		return tok, err
	}
	return tok, nil
}

func (c *Client) RevokeSession() error {
	return c.store.DeleteToken(c.issuer, c.clientID)
}

func login(issuer string, clientID string) (*session.Token, error) {
	flow, err := authflow.New(issuer, clientID)
	if err != nil {
		return nil, err
	}

	// TODO: Automatically open the browser at this URL or prompt the user.
	// TODO: handle non-interactive login flows.
	fmt.Printf("Press Enter to open your browser and login...")
	fmt.Scanln()

	err = flow.OpenBrowser()
	if err != nil {
		// Instead, should we print the URL and let the user open it themselves?
		return nil, err
	}

	/////////
	// TODO: technically we should start the callback server before opening the browser
	srv := callbackserver.New()
	err = srv.Listen()
	if err != nil {
		return nil, err
	}

	go func() {
		_ = srv.Start() // TODO: handle errors
	}()
	defer func() {
		_ = srv.Stop()
	}()
	resp := srv.WaitForResponse()
	// TODO: check state

	if resp.Error != "" {
		return nil, fmt.Errorf("error: %s", resp.Error)
	}
	code := resp.Code
	/////////

	// TODO: validate the state param that we'll receive in the redirect
	tok, err := flow.Exchange(code)
	if err != nil {
		return nil, err
	}

	return tok, nil
}
