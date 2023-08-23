// Code in this file taken from https://go-review.googlesource.com/c/oauth2/+/463979
// under a BSD license.
//
// The standard go library for oauth2 does not support PCKE, but a proposal
// has been accepted to add it in a future version of go. We copy the code from
// the proposal so that it's easy to switch to the standard library when it's
// available.
//
// Issue: https://github.com/golang/oauth2/issues/603
// Proposal: https://github.com/golang/go/issues/59835

package pkce

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/oauth2"
)

const (
	codeChallengeKey       = "code_challenge"
	codeChallengeMethodKey = "code_challenge_method"
	codeVerifierKey        = "code_verifier"
)

// GenerateVerifier generates a PKCE code verifier with 32 octets of randomness.
// This follows recommendations in RFC 7636.
//
// A fresh verifier should be generated for each authorization.
// S256ChallengeOption(verifier) should then be passed to Config.AuthCodeURL and
// VerifierOption(verifier) to Config.Exchange.
func GenerateVerifier() string {
	// "RECOMMENDED that the output of a suitable random number generator be
	// used to create a 32-octet sequence.  The octet sequence is then
	// base64url-encoded to produce a 43-octet URL-safe string to use as the
	// code verifier."
	// https://datatracker.ietf.org/doc/html/rfc7636#section-4.1
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

// VerifierOption describes a PKCE code verifier. It should be
// passed to Config.Exchange only.
func VerifierOption(verifier string) oauth2.AuthCodeOption {
	return oauth2.SetAuthURLParam(codeVerifierKey, verifier)
}

// S256ChallengeFromVerifier returns a PKCE code challenge derived from verifier with method S256.
//
// Prefer to use S256ChallengeOption where possible.
func S256ChallengeFromVerifier(verifier string) string {
	sha := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sha[:])
}

// S256ChallengeOption derives a PKCE code challenge derived from verifier with method S256.
// It should be passed to Config.AuthCodeURL only.
func S256ChallengeOption(verifier string) []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam(codeChallengeMethodKey, "S256"),
		oauth2.SetAuthURLParam(codeChallengeKey, S256ChallengeFromVerifier(verifier)),
	}
}
