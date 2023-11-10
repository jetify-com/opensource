package jetauth

import (
	"context"
	"fmt"
	"time"

	"go.jetpack.io/pkg/sandbox/auth"
	"go.jetpack.io/pkg/sandbox/auth/session"
	"go.jetpack.io/pkg/sandbox/jetcloud"
)

func NewClient(issuer, clientID string) (*auth.Client, error) {
	return auth.NewClient(
		issuer,
		clientID,
		getShortTermAccessToken,
		getShortTermAccessToken,
	)
}

func getShortTermAccessToken(
	ctx context.Context,
	tok *session.Token,
) (*session.Token, error) {
	accessToken, err := jetcloud.GetAccessToken(ctx, tok)
	if err != nil {
		fmt.Println(err)
		// We set the current set of tokens to be expired immediately.
		// This is a hack to force a refresh. Even though we have a new id token,
		// we failed to get a valid access token. This is likely because the
		// user doesn't have a valid plan.
		// Returning a nil token would prevent any new data from being written to
		// token store, but that means our existing refresh token would no longer
		// be valid because it is only usable once.
		tok.Expiry = time.Now()
		tok.AccessToken = ""
	} else {
		tok.AccessToken = accessToken
	}
	return tok, nil
}
