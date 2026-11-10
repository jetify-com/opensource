package tokenstore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"golang.org/x/oauth2"

	"go.jetify.com/pkg/auth/session"
)

const (
	testIssuer   = "https://auth.example.com/"
	testClientID = "client-123"
)

func newToken(t *testing.T, subject, accessToken string) *session.Token {
	t.Helper()
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       []byte("0123456789abcdef0123456789abcdef"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := jwt.Signed(signer).Claims(jwt.Claims{Subject: subject}).Serialize()
	if err != nil {
		t.Fatal(err)
	}
	return &session.Token{
		Token:   oauth2.Token{AccessToken: accessToken},
		IDToken: raw,
	}
}

func newStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(filepath.Join(t.TempDir(), "tokens"))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func accessTokens(toks []*session.Token) []string {
	out := make([]string, len(toks))
	for i, tok := range toks {
		out[i] = tok.AccessToken
	}
	return out
}

func assertAccessTokens(t *testing.T, got []*session.Token, want ...string) {
	t.Helper()
	gotStrs := accessTokens(got)
	if len(gotStrs) != len(want) {
		t.Fatalf("tokens = %v, want %v", gotStrs, want)
	}
	for i := range want {
		if gotStrs[i] != want[i] {
			t.Fatalf("tokens = %v, want %v", gotStrs, want)
		}
	}
}

func TestNewCreatesPrivateDir(t *testing.T) {
	root := filepath.Join(t.TempDir(), "a", "b")
	if _, err := New(root); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal("root is not a directory")
	}
	if perm := info.Mode().Perm(); perm&0o077 != 0 {
		t.Errorf("root dir perms = %v, want no group/other access", perm)
	}
}

func TestReadTokensEmpty(t *testing.T) {
	s := newStore(t)
	toks, err := s.ReadTokens(testIssuer, testClientID)
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 0 {
		t.Errorf("ReadTokens() = %v, want none", toks)
	}
}

func TestWriteTokenNil(t *testing.T) {
	s := newStore(t)
	if err := s.WriteToken(testIssuer, testClientID, nil, true); err == nil {
		t.Error("WriteToken(nil) error = nil, want error")
	}
}

func TestWriteAndReadDefaultToken(t *testing.T) {
	s := newStore(t)
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a1"), true); err != nil {
		t.Fatal(err)
	}
	toks, err := s.ReadTokens(testIssuer, testClientID)
	if err != nil {
		t.Fatal(err)
	}
	assertAccessTokens(t, toks, "a1")
	if got := toks[0].IDClaims(); got == nil || got.Subject != "alice" {
		t.Errorf("IDClaims() after round trip = %+v, want subject alice", got)
	}
}

func TestMakeDefaultOrdering(t *testing.T) {
	s := newStore(t)
	write := func(sub, access string) {
		t.Helper()
		if err := s.WriteToken(testIssuer, testClientID, newToken(t, sub, access), true); err != nil {
			t.Fatal(err)
		}
	}

	write("alice", "a1")
	write("bob", "b1")
	toks, _ := s.ReadTokens(testIssuer, testClientID)
	assertAccessTokens(t, toks, "b1", "a1")

	// Re-adding alice as default moves her to the front and replaces the
	// stale token rather than duplicating it.
	write("alice", "a2")
	toks, _ = s.ReadTokens(testIssuer, testClientID)
	assertAccessTokens(t, toks, "a2", "b1")
}

func TestReplaceTokenKeepsOrder(t *testing.T) {
	s := newStore(t)
	for _, tok := range []*session.Token{newToken(t, "alice", "a1"), newToken(t, "bob", "b1")} {
		if err := s.WriteToken(testIssuer, testClientID, tok, true); err != nil {
			t.Fatal(err)
		}
	}

	// bob is the default; refreshing alice's token must not change that.
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a2"), false); err != nil {
		t.Fatal(err)
	}
	toks, _ := s.ReadTokens(testIssuer, testClientID)
	assertAccessTokens(t, toks, "b1", "a2")
}

func TestReplaceTokenUnknownSubjectIsDropped(t *testing.T) {
	s := newStore(t)
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a1"), true); err != nil {
		t.Fatal(err)
	}
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "carol", "c1"), false); err != nil {
		t.Fatal(err)
	}
	toks, _ := s.ReadTokens(testIssuer, testClientID)
	assertAccessTokens(t, toks, "a1")
}

func TestIssuersAndClientsAreIsolated(t *testing.T) {
	s := newStore(t)
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a1"), true); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct{ issuer, client string }{
		{"https://other.example.com", testClientID},
		{testIssuer, "other-client"},
	} {
		toks, err := s.ReadTokens(tc.issuer, tc.client)
		if err != nil {
			t.Fatal(err)
		}
		if len(toks) != 0 {
			t.Errorf("ReadTokens(%q, %q) = %v, want none", tc.issuer, tc.client, accessTokens(toks))
		}
	}
}

func TestIssuerNormalization(t *testing.T) {
	s := newStore(t)
	if err := s.WriteToken("https://auth.example.com/", testClientID, newToken(t, "alice", "a1"), true); err != nil {
		t.Fatal(err)
	}
	// Scheme and trailing slash differences map to the same file.
	for _, issuer := range []string{"https://auth.example.com", "http://auth.example.com/", "auth.example.com"} {
		toks, err := s.ReadTokens(issuer, testClientID)
		if err != nil {
			t.Fatal(err)
		}
		assertAccessTokens(t, toks, "a1")
	}
}

func TestDeleteToken(t *testing.T) {
	s := newStore(t)

	// Deleting when nothing exists is a no-op.
	if err := s.DeleteToken(testIssuer, testClientID); err != nil {
		t.Errorf("DeleteToken() on missing file error = %v", err)
	}

	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a1"), true); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteToken(testIssuer, testClientID); err != nil {
		t.Fatal(err)
	}
	toks, err := s.ReadTokens(testIssuer, testClientID)
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 0 {
		t.Errorf("ReadTokens() after delete = %v, want none", accessTokens(toks))
	}
}

func TestReadTokensCorruptFile(t *testing.T) {
	s := newStore(t)
	path := s.path(testIssuer, testClientID)
	if err := ensureDir(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ReadTokens(testIssuer, testClientID); err == nil {
		t.Error("ReadTokens() error = nil, want JSON error")
	}
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a1"), true); err == nil {
		t.Error("WriteToken() error = nil, want error reading corrupt file")
	}
}

func TestStoredFileFormat(t *testing.T) {
	s := newStore(t)
	if err := s.WriteToken(testIssuer, testClientID, newToken(t, "alice", "a1"), true); err != nil {
		t.Fatal(err)
	}
	var data storeData
	if err := readJSONFile(s.path(testIssuer, testClientID), &data); err != nil {
		t.Fatal(err)
	}
	if data.Version != storeDataVersion {
		t.Errorf("Version = %q, want %q", data.Version, storeDataVersion)
	}
	if len(data.Tokens) != 1 || data.Tokens[0].IDToken == "" {
		t.Errorf("stored tokens = %+v, want one token with an ID token", data.Tokens)
	}
}

func TestIssuerSlug(t *testing.T) {
	tests := map[string]string{
		"https://auth.example.com/": "auth-example-com",
		"http://auth.example.com":   "auth-example-com",
		"auth.example.com":          "auth-example-com",
		"https://Auth.Example.com/": "auth-example-com",
	}
	for in, want := range tests {
		if got := issuerSlug(in); got != want {
			t.Errorf("issuerSlug(%q) = %q, want %q", in, got, want)
		}
	}
}
