package auth

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.jetify.com/pkg/auth/session"
	"golang.org/x/oauth2"

	"go.jetify.com/pkg/auth/internal/authflow"
	"go.jetify.com/pkg/auth/internal/callbackserver"
	"go.jetify.com/pkg/auth/internal/tokenstore"
)

var ErrNotLoggedIn = fmt.Errorf("not logged in")

type Client struct {
	audience        []string
	issuer          string
	clientID        string
	store           *tokenstore.Store
	scopes          []string
	successRedirect string
}

func NewClient(
	issuer, clientID string,
	scopes []string,
	successRedirect string,
	audience []string,
) (*Client, error) {
	store, err := tokenstore.New(storeDir())
	if err != nil {
		return nil, err
	}

	return &Client{
		audience:        audience,
		issuer:          issuer,
		clientID:        clientID,
		scopes:          scopes,
		store:           store,
		successRedirect: successRedirect,
	}, nil
}

func storeDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "~/.cache"
	}
	return filepath.Join(cacheDir, "jetify", "auth")
}

func (c *Client) LoginFlow() (*session.Token, error) {
	tok, err := c.login()
	if err != nil {
		return nil, err
	}

	_ = c.store.WriteToken(c.issuer, c.clientID, tok, true /*makeDefault*/)
	return tok, nil
}

func (c *Client) LoginFlowIfNeeded(ctx context.Context) (*session.Token, error) {
	tok, err := c.GetSession(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "You are not logged in.")
		tok, err = c.LoginFlow()
	}
	return tok, err
}

// LoginFlowIfNeededForOrg returns the current valid session token for a given
// org and prompts to log in if needed.
// Note: I'm not sure this is best API. Currently evolving
func (c *Client) LoginFlowIfNeededForOrg(
	ctx context.Context,
	orgID string,
) (*session.Token, error) {
	sessions, err := c.GetSessions()
	if err != nil {
		return nil, err
	}
	for _, session := range sessions {
		if session.Peek().IDClaims().OrgID == orgID {
			return session.Token(ctx)
		}
	}
	if len(sessions) > 0 {
		fmt.Fprintln(os.Stderr, "You are not logged in to organization that owns this project. Please log in.")
	} else {
		fmt.Fprintln(os.Stderr, "You are not logged in.")
	}
	return c.LoginFlow()
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
	sources, err := c.GetSessions()
	if err != nil {
		return nil, err
	}

	if len(sources) == 0 || sources[0].Peek() == nil {
		return nil, ErrNotLoggedIn
	}

	return sources[0].Token(ctx)
}

// GetSessions returns all session tokens as refreshableTokenSource. This means
// that they may not be valid or even refreshable. Callers can use Token() to
// refresh the token if needed. Even if Token() fails to refresh, it still
// returns the token, so callers can use the expired data if they want.
// This allows callers to list all sessions, even if they are expired.
// Callers can use Peek() to inspect token data without refreshing.
func (c *Client) GetSessions() ([]refreshableTokenSource, error) {
	tokens, err := c.store.ReadTokens(c.issuer, c.clientID)
	if err != nil {
		return nil, err
	}

	refreshableTokens := []refreshableTokenSource{}
	for _, tok := range tokens {
		refreshableTokens = append(refreshableTokens, refreshableTokenSource{
			client: c,
			token:  tok,
		})
	}

	return refreshableTokens, nil
}

func (c *Client) AddSession(tok *session.Token) error {
	return c.store.WriteToken(c.issuer, c.clientID, tok, true /*makeDefault*/)
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
		Scopes:   c.scopes,
	}

	// Refresh logic:
	tokenSource := conf.TokenSource(ctx, &tok.Token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return tok, err
	}

	if newToken.AccessToken != tok.AccessToken {
		tok.Token = *newToken
		tok.IDToken = newToken.Extra("id_token").(string)
		err = c.store.WriteToken(c.issuer, c.clientID, tok, false /*makeDefault*/)
		if err != nil {
			return tok, err
		}
	}

	return tok, nil
}

func (c *Client) RevokeSession() error {
	return c.store.DeleteToken(c.issuer, c.clientID)
}

func (c *Client) login() (*session.Token, error) {
	flow, err := authflow.New(c.issuer, c.clientID, c.scopes, c.audience)
	if err != nil {
		return nil, err
	}

	// TODO: Automatically open the browser at this URL or prompt the user.
	// TODO: handle non-interactive login flows.
	fmt.Printf("Press Enter to open your browser and login...")
	_, err = fmt.Scanln()
	if err != nil {
		return nil, err
	}

	err = flow.OpenBrowser()
	if err != nil {
		// Instead, should we print the URL and let the user open it themselves?
		return nil, err
	}

	/////////
	// TODO: technically we should start the callback server before opening the browser
	srv := callbackserver.New(c.successRedirect)
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
