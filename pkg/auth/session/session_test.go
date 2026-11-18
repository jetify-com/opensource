package session

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"golang.org/x/oauth2"
)

func signIDToken(t *testing.T, alg jose.SignatureAlgorithm, key any, claims IDClaims) string {
	t.Helper()
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: alg, Key: key}, nil)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := jwt.Signed(signer).Claims(claims).Serialize()
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func hsKey() []byte { return []byte("0123456789abcdef0123456789abcdef") }

func TestIDClaims(t *testing.T) {
	want := IDClaims{
		Issuer:        "https://auth.example.com/",
		Subject:       "user-123",
		Audience:      jwt.Audience{"client-id"},
		Email:         "someone@example.com",
		EmailVerified: true,
		Name:          "Some One",
		OrgID:         "org_abc",
	}
	tok := &Token{IDToken: signIDToken(t, jose.HS256, hsKey(), want)}

	got := tok.IDClaims()
	if got == nil {
		t.Fatal("IDClaims() = nil")
	}
	if got.Issuer != want.Issuer || got.Subject != want.Subject ||
		got.Email != want.Email || !got.EmailVerified ||
		got.Name != want.Name || got.OrgID != want.OrgID {
		t.Errorf("IDClaims() = %+v, want %+v", got, want)
	}
	if !got.Audience.Contains("client-id") {
		t.Errorf("Audience = %v, want to contain client-id", got.Audience)
	}
}

func TestIDClaimsES256(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tok := &Token{IDToken: signIDToken(t, jose.ES256, key, IDClaims{Subject: "es"})}
	if c := tok.IDClaims(); c == nil || c.Subject != "es" {
		t.Errorf("IDClaims() = %+v, want subject es", c)
	}
}

func TestIDClaimsInvalid(t *testing.T) {
	tests := map[string]string{
		"empty":        "",
		"garbage":      "not-a-jwt",
		"bad segments": "a.b.c",
		// Header {"alg":"none"} with an empty signature. Unsigned tokens
		// must be rejected.
		"alg none": "eyJhbGciOiJub25lIn0.eyJzdWIiOiJ4In0.",
	}
	for name, raw := range tests {
		t.Run(name, func(t *testing.T) {
			tok := &Token{IDToken: raw}
			if c := tok.IDClaims(); c != nil {
				t.Errorf("IDClaims() = %+v, want nil", c)
			}
		})
	}
}

func TestIDClaimsCached(t *testing.T) {
	cached := &IDClaims{Subject: "cached"}
	tok := &Token{IDToken: "", idClaims: cached}
	if got := tok.IDClaims(); got != cached {
		t.Errorf("IDClaims() = %v, want cached value", got)
	}
}

func TestFromOauth2Nil(t *testing.T) {
	tok, err := FromOauth2(nil)
	if err != nil || tok != nil {
		t.Errorf("FromOauth2(nil) = %v, %v; want nil, nil", tok, err)
	}
}

func TestFromOauth2WithoutIDToken(t *testing.T) {
	otok := &oauth2.Token{AccessToken: "access", RefreshToken: "refresh"}
	tok, err := FromOauth2(otok)
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "access" || tok.RefreshToken != "refresh" {
		t.Errorf("embedded oauth2 token not copied: %+v", tok.Token)
	}
	if tok.IDToken != "" || tok.Issuer != "" {
		t.Errorf("IDToken = %q, Issuer = %q; want both empty", tok.IDToken, tok.Issuer)
	}
	if tok.IDClaims() != nil {
		t.Error("IDClaims() should be nil without an ID token")
	}
}

func TestFromOauth2WithIDToken(t *testing.T) {
	raw := signIDToken(t, jose.HS256, hsKey(), IDClaims{
		Issuer:  "https://issuer.example.com",
		Subject: "sub",
	})
	otok := (&oauth2.Token{AccessToken: "access"}).WithExtra(map[string]any{
		"id_token": raw,
	})

	tok, err := FromOauth2(otok)
	if err != nil {
		t.Fatal(err)
	}
	if tok.IDToken != raw {
		t.Errorf("IDToken = %q, want the raw id_token extra", tok.IDToken)
	}
	if tok.Issuer != "https://issuer.example.com" {
		t.Errorf("Issuer = %q, want issuer from claims", tok.Issuer)
	}
}

func TestFromOauth2WithUnparseableIDToken(t *testing.T) {
	otok := (&oauth2.Token{}).WithExtra(map[string]any{"id_token": "junk"})
	tok, err := FromOauth2(otok)
	if err != nil {
		t.Fatal(err)
	}
	if tok.IDToken != "junk" {
		t.Errorf("IDToken = %q, want junk", tok.IDToken)
	}
	if tok.Issuer != "" {
		t.Errorf("Issuer = %q, want empty", tok.Issuer)
	}
}

func TestVerify(t *testing.T) {
	if err := (&Token{}).Verify(); err != nil {
		t.Errorf("Verify() = %v", err)
	}
	if err := (&Token{Keys: []jose.JSONWebKey{}}).Verify(); err != nil {
		t.Errorf("Verify() = %v", err)
	}
}
