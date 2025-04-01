package httpmock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// ReplayServer is a server that can be used to test HTTP interactions by
// recording and replaying real HTTP requests.
//
// When a test is first run, the ReplayServer will record the interactions
// and save them to a "cassette" file.
//
// On subsequent runs, the ReplayServer will replay the interactions from the
// cassette file, allowing for consistent tests.
type ReplayServer struct {
	server           *httptest.Server
	rec              *recorder.Recorder
	t                T
	realURL          *url.URL
	cassetteName     string
	usedInteractions atomic.Int32
}

// ReplayConfig is the configuration for a ReplayServer.
type ReplayConfig struct {
	// Host is the address of the real endpoints we're proxying requests to,
	// e.g. "https://api.example.com".
	Host string
	// Cassette is the name of the cassette file to use for recording and replaying.
	// If the cassette does not exist, it will be created.
	// Do not include the ".yaml" extension, it will be added automatically.
	Cassette string
}

// NewReplayServer creates a new ReplayServer.
func NewReplayServer(tester T, config ReplayConfig) (*ReplayServer, error) {
	tester.Helper()

	// Parse and validate the real endpoint we'll proxy to in record mode
	realURL, err := url.Parse(config.Host)
	if err != nil {
		return nil, fmt.Errorf("invalid Host: %v", err)
	}
	if realURL.Scheme == "" || realURL.Host == "" {
		return nil, fmt.Errorf("invalid Host: URL must have scheme and host")
	}

	// Create ReplayServer first
	replayServer := &ReplayServer{
		realURL:      realURL,
		t:            tester,
		cassetteName: config.Cassette,
	}

	rec, err := recorder.New(
		config.Cassette,
		recorder.WithMode(recorder.ModeRecordOnce),
		recorder.WithHook(removeIgnored, recorder.AfterCaptureHook),
		recorder.WithMatcher(func(request *http.Request, cassetteRequest cassette.Request) bool {
			return requireCassetteRequest(tester, cassetteRequest, request)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create recorder: %v", err)
	}

	replayServer.rec = rec
	replayServer.server = httptest.NewServer(http.HandlerFunc(replayServer.handler))

	return replayServer, nil
}

// handler processes incoming requests by either replaying them from the cassette
// or proxying them to the real endpoint and recording the result.
func (rs *ReplayServer) handler(w http.ResponseWriter, req *http.Request) {
	// Create a copy of the request for matching
	matchReq := *req

	// We'll rewrite the request's scheme and host to point at the real endpoint
	matchReq.URL.Scheme = rs.realURL.Scheme
	matchReq.URL.Host = rs.realURL.Host
	matchReq.Host = rs.realURL.Host

	// Now let the recorder handle it (replay or forward+record)
	resp, roundTripErr := rs.rec.RoundTrip(&matchReq)
	if roundTripErr != nil {
		// Log the error details for debugging
		rs.t.Errorf("Replay server error: %v", roundTripErr)
		// Enhance error reporting for request matching failures
		http.Error(w, fmt.Sprintf("request did not match cassette: %v", roundTripErr), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Increment the used interactions counter
	rs.usedInteractions.Add(1)

	// Copy response status, headers, body
	for k, vals := range resp.Header {
		for _, val := range vals {
			w.Header().Add(k, val)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, copyErr := io.Copy(w, resp.Body); copyErr != nil {
		rs.t.Errorf("Error copying response body: %v", copyErr)
		http.Error(w, fmt.Sprintf("error copying response: %v", copyErr), http.StatusInternalServerError)
		require.NoError(rs.t, copyErr)
	}
}

// formatRequestBody reads and returns the request body as a string, restoring it afterward
func formatRequestBody(req *http.Request) string {
	if req.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Sprintf("error reading body: %v", err)
	}
	// Restore the body for later use
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes)
}

// Close stops the ReplayServer and verifies that all recorded interactions
// were used.
func (rs *ReplayServer) Close() error {
	// Stop the HTTP server first
	rs.server.Close()

	// Verify that all recorded interactions were used
	if err := rs.rec.Stop(); err != nil {
		return err
	}

	// Get the cassette to check if all interactions were used
	cassette, err := cassette.Load(rs.cassetteName)
	if err != nil {
		return fmt.Errorf("failed to load cassette for verification: %v", err)
	}

	// Check if any interactions were not used
	usedCount := rs.usedInteractions.Load()
	if usedCount < int32(len(cassette.Interactions)) {
		nextUnused := cassette.Interactions[usedCount]
		return fmt.Errorf("expected %d requests, received %d. Next expected: [%s %s]",
			len(cassette.Interactions), usedCount, nextUnused.Request.Method, nextUnused.Request.URL)
	}

	// Check if we received more requests than expected
	if usedCount > int32(len(cassette.Interactions)) {
		return fmt.Errorf("expected %d requests, received %d: too many requests made",
			len(cassette.Interactions), usedCount)
	}

	return nil
}

// URL returns the base URL of the replay server.
func (rs *ReplayServer) URL() string {
	return rs.server.URL
}

var ignoredHeaders = []string{
	// Sensitive headers:
	"Authorization",
	"Proxy-Authorization",
	"WWW-Authenticate",
	"Proxy-Authenticate",
	"Cookie",
	"Set-Cookie",
	"X-API-Key",
	"X-Auth-Token",
	"X-CSRF-Token",
	"X-Requested-With",
	"X-Forwarded-For",
	"X-Real-IP",
	"X-Client-IP",

	// Other headers:
	"User-Agent",
	"Accept",
	"Accept-Encoding",
	"Accept-Language",
	"Connection",
	"Content-Length",
}

// removeIgnored is a recorder hook that removes ignored headers and request fields
// before they are saved to the cassette.
func removeIgnored(i *cassette.Interaction) error {
	// Remove request headers
	for _, header := range ignoredHeaders {
		i.Request.Headers.Del(header)
	}
	// Remove response headers
	for _, header := range ignoredHeaders {
		i.Response.Headers.Del(header)
	}
	// Remove remote address from request
	i.Request.RemoteAddr = ""
	return nil
}
