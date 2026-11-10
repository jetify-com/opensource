package tokenstore

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"go.jetify.com/pkg/auth/session"
	"golang.org/x/oauth2"
)

const (
	testIssuer   = "https://auth.example.com"
	testClientID = "test-client"
)

// newToken returns a token whose ID token is an HS256-signed JWT with the given
// subject. accessToken lets tests tell tokens with the same subject apart.
func newToken(t *testing.T, sub, accessToken string) *session.Token {
	t.Helper()
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte("0123456789abcdef0123456789abcdef")},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		t.Fatal(err)
	}
	idToken, err := jwt.Signed(signer).Claims(jwt.Claims{Subject: sub, Issuer: testIssuer}).Serialize()
	if err != nil {
		t.Fatal(err)
	}
	return &session.Token{
		Token:   oauth2.Token{AccessToken: accessToken},
		IDToken: idToken,
	}
}

func newStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// writeRaw writes tokens directly to the store file, bypassing WriteToken's
// validation, to simulate a store that already holds bad tokens.
func writeRaw(t *testing.T, s *Store, tokens ...*session.Token) {
	t.Helper()
	err := writeJSONFile(s.path(testIssuer, testClientID), storeData{
		Version: storeDataVersion,
		Tokens:  tokens,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func accessTokens(t *testing.T, s *Store) []string {
	t.Helper()
	toks, err := s.ReadTokens(testIssuer, testClientID)
	if err != nil {
		t.Fatal(err)
	}
	out := []string{}
	for _, tok := range toks {
		out = append(out, tok.AccessToken)
	}
	return out
}

func assertAccessTokens(t *testing.T, s *Store, want ...string) {
	t.Helper()
	got := accessTokens(t, s)
	gotJSON, _ := json.Marshal(got)
	wantJSON, _ := json.Marshal(want)
	if string(gotJSON) != string(wantJSON) {
		t.Errorf("stored tokens = %s, want %s", gotJSON, wantJSON)
	}
}

func TestWriteTokenDefaultAndReplace(t *testing.T) {
	s := newStore(t)
	for _, tok := range []*session.Token{newToken(t, "alice", "a1"), newToken(t, "bob", "b1")} {
		if err := s.WriteToken(testIssuer, testClientID, tok, true); err != nil {
			t.Fatal(err)
		}
	}
	assertAccessTokens(t, s, "b1", "a1")

	// Making alice default again moves her to the front and replaces her token.
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a2"), true); err != nil {
		t.Fatal(err)
	}
	assertAccessTokens(t, s, "a2", "b1")

	// Replacing bob keeps his position.
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "bob", "b2"), false); err != nil {
		t.Fatal(err)
	}
	assertAccessTokens(t, s, "a2", "b2")
}

func TestWriteTokenRejectsTokenWithoutClaims(t *testing.T) {
	cases := map[string]*session.Token{
		"no id token":   {Token: oauth2.Token{AccessToken: "x"}},
		"malformed":     {Token: oauth2.Token{AccessToken: "x"}, IDToken: "not-a-jwt"},
		"empty subject": newToken(t, "", "x"),
	}
	for name, tok := range cases {
		for _, makeDefault := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s/makeDefault=%v", name, makeDefault), func(t *testing.T) {
				s := newStore(t)
				writeRaw(t, s, newToken(t, "alice", "a1"))

				if err := s.WriteToken(testIssuer, testClientID, tok, makeDefault); err == nil {
					t.Fatal("WriteToken succeeded, want error")
				}
				// The store is left untouched.
				assertAccessTokens(t, s, "a1")
			})
		}
	}
}

func TestWriteTokenNoStoreFileRejectsTokenWithoutClaims(t *testing.T) {
	s := newStore(t)
	tok := &session.Token{Token: oauth2.Token{AccessToken: "x"}}
	if err := s.WriteToken(testIssuer, testClientID, tok, true); err == nil {
		t.Fatal("WriteToken succeeded, want error")
	}
	if _, err := os.Stat(s.path(testIssuer, testClientID)); !os.IsNotExist(err) {
		t.Errorf("store file exists after rejected write: %v", err)
	}
}

func TestWriteTokenDropsStoredTokensWithoutClaims(t *testing.T) {
	stored := func(t *testing.T) []*session.Token {
		return []*session.Token{
			newToken(t, "alice", "a1"),
			{Token: oauth2.Token{AccessToken: "no-id"}},
			newToken(t, "bob", "b1"),
			{Token: oauth2.Token{AccessToken: "bad-id"}, IDToken: "not-a-jwt"},
		}
	}

	t.Run("make default", func(t *testing.T) {
		s := newStore(t)
		writeRaw(t, s, stored(t)...)
		if err := s.WriteToken(testIssuer, testClientID, newToken(t, "bob", "b2"), true); err != nil {
			t.Fatal(err)
		}
		assertAccessTokens(t, s, "b2", "a1")
	})

	t.Run("replace", func(t *testing.T) {
		s := newStore(t)
		writeRaw(t, s, stored(t)...)
		if err := s.WriteToken(testIssuer, testClientID, newToken(t, "bob", "b2"), false); err != nil {
			t.Fatal(err)
		}
		assertAccessTokens(t, s, "a1", "b2")
	})
}
