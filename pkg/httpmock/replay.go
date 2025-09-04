package httpmock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
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
	// Mode determines whether to use cassettes (unit testing) or always hit the real API (integration)
	// If not specified, defaults to the value of INTEGRATION environment variable
	//
	// Possible values are "unit" and "integration".
	Mode TestMode
	// AdditionalIgnoredHeaders specifies additional headers to ignore during recording and replay.
	// These headers will be added to the default set of ignored headers.
	// Header names are case-insensitive.
	AdditionalIgnoredHeaders []string
}

// TestMode represents how the ReplayServer should handle HTTP interactions
type TestMode string

const (
	// ModeUnitTest uses cassettes for replay (default behavior)
	ModeUnitTest TestMode = "unit"
	// ModeIntegrationTest always forwards requests to the real API
	ModeIntegrationTest TestMode = "integration"
)

// NewReplayServer creates a new ReplayServer.
func NewReplayServer(tester T, config ReplayConfig) *ReplayServer {
	tester.Helper()

	// Parse and validate the real endpoint we'll proxy to in record mode
	realURL, err := url.Parse(config.Host)
	require.NoError(tester, err, "invalid Host URL: %s", config.Host)
	require.True(tester, realURL.Scheme != "" && realURL.Host != "", "Host URL must have scheme and host, got: %s", config.Host)

	// Create ReplayServer first
	replayServer := &ReplayServer{
		realURL:      realURL,
		t:            tester,
		cassetteName: config.Cassette,
	}

	// Create the recorder
	rec, err := createRecorder(tester, config)
	require.NoError(tester, err, "failed to create recorder")

	replayServer.rec = rec
	replayServer.server = httptest.NewServer(http.HandlerFunc(replayServer.handler))

	return replayServer
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
	defer func() { _ = resp.Body.Close() }()

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

// Assert verifies that all recorded interactions were used. This can be called
// separately from Close() to perform assertions without stopping the server.
func (rs *ReplayServer) assert() error {
	// Get the cassette to check if all interactions were used
	cassette, err := cassette.Load(rs.cassetteName)
	if err != nil {
		return fmt.Errorf("failed to load cassette for verification: %w", err)
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

// Close stops the ReplayServer and verifies that all recorded interactions
// were used.
func (rs *ReplayServer) Close() {
	// Stop the HTTP server first
	rs.server.Close()

	// Verify that all recorded interactions were used
	err := rs.rec.Stop()
	require.NoError(rs.t, err, "failed to stop recorder")

	// Verify that all recorded interactions were used
	err = rs.assert()
	require.NoError(rs.t, err, "not all recorded interactions were used")
}

// URL returns the base URL of the replay server.
func (rs *ReplayServer) URL() string {
	return rs.server.URL
}

var ignoredHeaders = []string{
	// Sensitive headers:
	"_csrf",
	"_csrf_token",
	"_session",
	"_xsrf",
	"Api-Key",
	"Apikey",
	"Auth",
	"Authorization",
	"Cookie",
	"Credentials",
	"Csrf",
	"Csrf-Token",
	"Csrftoken",
	"Ip-Address",
	"Passwd",
	"Password",
	"Private-Key",
	"Privatekey",
	"Proxy-Authorization",
	"Proxy-Authenticate",
	"Remote-Addr",
	"Secret",
	"Session",
	"Sessionid",
	"Set-Cookie",
	"Token",
	"User-Session",
	"WWW-Authenticate",
	"X-Api-Key",
	"X-API-Key",
	"X-Auth-Token",
	"X-Client-IP",
	"X-CSRF-Token",
	"X-Csrftoken",
	"X-Forwarded-For",
	"X-Real-IP",
	"X-Real-Ip",
	"X-Requested-With",
	"XSRF-TOKEN",

	// Other headers:
	"User-Agent",
	"Accept",
	"Accept-Encoding",
	"Accept-Language",
	"Connection",
	"Content-Length",
}

// createRecorder creates a new recorder with the given configuration
func createRecorder(tester T, config ReplayConfig) (*recorder.Recorder, error) {
	// Determine the recorder mode based on config and environment
	recorderMode := recorder.ModeRecordOnce
	if config.Mode == "" {
		// If mode not specified, use environment variable setting
		if val, err := strconv.ParseBool(os.Getenv("INTEGRATION")); err == nil && val {
			config.Mode = ModeIntegrationTest
		} else {
			config.Mode = ModeUnitTest
		}
	}
	if config.Mode == ModeIntegrationTest {
		recorderMode = recorder.ModePassthrough
	}

	allIgnoredHeaders := append(ignoredHeaders, config.AdditionalIgnoredHeaders...)

	rec, err := recorder.New(
		config.Cassette,
		recorder.WithMode(recorderMode),
		recorder.WithHook(func(i *cassette.Interaction) error {
			return removeIgnoredHeaders(i, allIgnoredHeaders)
		}, recorder.AfterCaptureHook),
		recorder.WithHook(decompressGzip, recorder.BeforeSaveHook),
		recorder.WithMatcher(func(request *http.Request, cassetteRequest cassette.Request) bool {
			return requireCassetteRequest(tester, cassetteRequest, request)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create recorder: %w", err)
	}
	return rec, nil
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

// removeIgnoredHeaders is a recorder hook that removes ignored headers and request fields
// before they are saved to the cassette.
func removeIgnoredHeaders(i *cassette.Interaction, headersToIgnore []string) error {
	// Remove request headers
	for _, header := range headersToIgnore {
		i.Request.Headers.Del(header)
	}
	// Remove response headers
	for _, header := range headersToIgnore {
		i.Response.Headers.Del(header)
	}
	// Remove remote address from request
	i.Request.RemoteAddr = ""
	return nil
}
