package auth

import (
	"context"

	"go.jetpack.io/pkg/auth/session"
)

// Implements something similar to oauth2.TokenSource, but with a session.Token
type refreshableToken struct {
	client *Client
	token  *session.Token
}

func (t *refreshableToken) Token(ctx context.Context) (*session.Token, error) {
	if !t.token.Valid() {
		return t.client.refresh(ctx, t.token)
	}
	return t.token, nil
}

func (t *refreshableToken) Peek() *session.Token {
	return t.token
}
