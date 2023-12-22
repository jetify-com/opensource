package auth

import (
	"context"

	"go.jetpack.io/pkg/auth/session"
)

// Implements something similar to oauth2.TokenSource, but with a session.Token
type refreshableTokenSource struct {
	client *Client
	token  *session.Token
}

func (t *refreshableTokenSource) Token(
	ctx context.Context,
) (*session.Token, error) {
	if !t.token.Valid() {
		return t.client.refresh(ctx, t.token)
	}
	return t.token, nil
}

func (t *refreshableTokenSource) Peek() *session.Token {
	return t.token
}
