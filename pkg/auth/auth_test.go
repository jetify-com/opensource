package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"go.jetify.com/pkg/auth/internal/tokenstore"
	"go.jetify.com/pkg/auth/session"
	"golang.org/x/oauth2"
)

func makeIDToken(t *testing.T, sub, name string) string {
	t.Helper()
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte("0123456789abcdef0123456789abcdef")},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := jwt.Signed(signer).Claims(map[string]any{"sub": sub, "name": name}).Serialize()
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// newTestIssuer starts a minimal OIDC issuer whose token endpoint returns
// tokenResp for refresh requests.
func newTestIssuer(t *testing.T, tokenResp map[string]any) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":                 srv.URL,
			"authorization_endpoint": srv.URL + "/authorize",
			"token_endpoint":         srv.URL + "/token",
			"jwks_uri":               srv.URL + "/jwks",
		})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tokenResp)
	})
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestRefreshWithoutIDTokenKeepsExistingIDToken(t *testing.T) {
	oldID := makeIDToken(t, "user-1", "old")
	srv := newTestIssuer(t, map[string]any{
		"access_token":  "new-access",
		"token_type":    "Bearer",
		"refresh_token": "new-refresh",
		"expires_in":    3600,
		// No id_token: allowed by OIDC for refresh responses.
	})

	store, err := tokenstore.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{issuer: srv.URL, clientID: "client", store: store}

	tok := &session.Token{
		Token: oauth2.Token{
			AccessToken:  "old-access",
			RefreshToken: "old-refresh",
			Expiry:       time.Now().Add(-time.Hour),
		},
		IDToken: oldID,
	}

	got, err := c.refresh(context.Background(), tok)
	if err != nil {
		t.Fatalf("refresh() error = %v", err)
	}
	if got.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want %q", got.AccessToken, "new-access")
	}
	if got.IDToken != oldID {
		t.Errorf("IDToken = %q, want existing %q", got.IDToken, oldID)
	}
}

func TestRefreshWithIDTokenReplacesIt(t *testing.T) {
	oldID := makeIDToken(t, "user-1", "old")
	newID := makeIDToken(t, "user-1", "new")
	srv := newTestIssuer(t, map[string]any{
		"access_token": "new-access",
		"token_type":   "Bearer",
		"expires_in":   3600,
		"id_token":     newID,
	})

	store, err := tokenstore.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{issuer: srv.URL, clientID: "client", store: store}

	tok := &session.Token{
		Token: oauth2.Token{
			AccessToken:  "old-access",
			RefreshToken: "old-refresh",
			Expiry:       time.Now().Add(-time.Hour),
		},
		IDToken: oldID,
	}

	got, err := c.refresh(context.Background(), tok)
	if err != nil {
		t.Fatalf("refresh() error = %v", err)
	}
	if got.IDToken != newID {
		t.Errorf("IDToken = %q, want %q", got.IDToken, newID)
	}
}
