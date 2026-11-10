package pkce

import (
	"encoding/base64"
	"net/url"
	"testing"

	"golang.org/x/oauth2"
)

func TestGenerateVerifier(t *testing.T) {
	v := GenerateVerifier()

	// 32 random bytes, base64url encoded without padding, is 43 chars.
	// RFC 7636 requires a verifier between 43 and 128 chars.
	if len(v) != 43 {
		t.Errorf("len(GenerateVerifier()) = %d, want 43", len(v))
	}
	raw, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		t.Fatalf("verifier %q is not raw base64url: %v", v, err)
	}
	if len(raw) != 32 {
		t.Errorf("verifier decodes to %d bytes, want 32", len(raw))
	}

	if GenerateVerifier() == v {
		t.Error("GenerateVerifier() returned the same verifier twice")
	}
}

func TestS256ChallengeFromVerifier(t *testing.T) {
	// Test vector from RFC 7636, Appendix B.
	const (
		verifier  = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		challenge = "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
	)
	if got := S256ChallengeFromVerifier(verifier); got != challenge {
		t.Errorf("S256ChallengeFromVerifier(%q) = %q, want %q", verifier, got, challenge)
	}
}

func authURLParams(t *testing.T, opts ...oauth2.AuthCodeOption) url.Values {
	t.Helper()
	cfg := &oauth2.Config{
		ClientID: "client",
		Endpoint: oauth2.Endpoint{AuthURL: "https://example.com/authorize"},
	}
	u, err := url.Parse(cfg.AuthCodeURL("state", opts...))
	if err != nil {
		t.Fatal(err)
	}
	return u.Query()
}

func TestS256ChallengeOption(t *testing.T) {
	const verifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	q := authURLParams(t, S256ChallengeOption(verifier)...)

	if got := q.Get(codeChallengeMethodKey); got != "S256" {
		t.Errorf("%s = %q, want S256", codeChallengeMethodKey, got)
	}
	if got, want := q.Get(codeChallengeKey), S256ChallengeFromVerifier(verifier); got != want {
		t.Errorf("%s = %q, want %q", codeChallengeKey, got, want)
	}
	if q.Has(codeVerifierKey) {
		t.Errorf("challenge option must not leak the verifier into the auth URL")
	}
}

func TestVerifierOption(t *testing.T) {
	const verifier = "my-verifier"
	q := authURLParams(t, VerifierOption(verifier))
	if got := q.Get(codeVerifierKey); got != verifier {
		t.Errorf("%s = %q, want %q", codeVerifierKey, got, verifier)
	}
}
