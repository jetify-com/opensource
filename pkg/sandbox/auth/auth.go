package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"go.jetpack.io/pkg/sandbox/auth/session"

	"go.jetpack.io/pkg/sandbox/auth/internal/authflow"
	"go.jetpack.io/pkg/sandbox/auth/internal/callbackserver"
	"go.jetpack.io/pkg/sandbox/auth/internal/tokenstore"
)

type Client struct {
	issuer   string
	clientID string
	store    *tokenstore.Store
}

func NewClient(issuer string, clientID string) (*Client, error) {
	store, err := tokenstore.New(storeDir())
	if err != nil {
		return nil, err
	}

	return &Client{
		issuer:   issuer,
		clientID: clientID,
		store:    store,
	}, nil
}

func storeDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "~/.cache"
	}
	return filepath.Join(cacheDir, "go.jetpack.io", "auth")
}

func (c *Client) LoginFlow() (*session.Token, error) {
	tok, err := login(c.issuer, c.clientID)
	if err != nil {
		return nil, err
	}

	_ = c.store.WriteToken(c.issuer, c.clientID, tok)
	return tok, nil
}

func (c *Client) LogoutFlow() error {
	// For now we just delete the token from the store.
	// But in the future we might want to revoke the token with the server, and do the oauth logout flow.
	return c.RevokeSession()
}

func (c *Client) GetSession() *session.Token {
	// TODO: automatically refresh token as needed
	return c.store.ReadToken(c.issuer, c.clientID)
}

func (c *Client) RefreshSession() *session.Token {
	panic("refresh session not implemented")
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
