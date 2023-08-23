package auth

import (
	"fmt"

	"go.jetpack.io/auth/internal/authflow"
	"go.jetpack.io/auth/internal/callbackserver"
	"golang.org/x/oauth2"
)

type Client struct {
	issuer   string
	clientID string
}

func NewClient(issuer string, clientID string) (*Client, error) {
	return &Client{
		issuer:   issuer,
		clientID: clientID,
	}, nil
}

func Login(issuer string, clientID string) (*oauth2.Token, error) {
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

	go srv.Start() // TODO: handle errors
	defer srv.Stop()
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
