package httpcacher

import "net/http"

// Implementation of an HTTP client with caching. Under the hoos we're using
// https://github.com/gregjones/httpcache which seems to be the 'best' one
// available.
//
// Note however that as of 2023, it is not actively maintained, it was originally
// developed 11 years ago, and was last updated 4 years ago. We should instead
// consider using an imlementation based on https://github.com/pquerna/cachecontrol
// Like in:
// + https://github.com/dadrus/heimdall/blob/main/internal/httpcache/round_tripper.go
// + https://github.com/darkweak/souin/blob/master/pkg/middleware/middleware.go
// It could even implement state-while-revalidate type of logic on the client
// side: https://developer.mozilla.org/en-US/docs/Web/API/Request/cache

var DefaultClient = NewClient(defaultCacheDir)

func NewClient(cacheDir string) *http.Client {
	return newTransport(cacheDir).Client()
}
