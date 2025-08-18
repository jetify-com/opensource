package httpmock

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// ReplayTransport implements http.RoundTripper for recording and replaying
// HTTP interactions across multiple hosts. Unlike ReplayServer which acts as
// a proxy, ReplayTransport can be used directly as the Transport for an
// http.Client, allowing your code to make requests to multiple different
// hosts while recording/replaying all interactions.
//
// When a test is first run, the ReplayTransport will record the interactions
// and save them to a "cassette" file.
//
// On subsequent runs, the ReplayTransport will replay the interactions from the
// cassette file, allowing for consistent tests.
type ReplayTransport struct {
	rec              *recorder.Recorder
	t                T
	cassetteName     string
	usedInteractions atomic.Int32
}

// NewReplayTransport creates a new ReplayTransport that can be used as an
// http.RoundTripper for recording and replaying HTTP interactions.
//
// Example usage:
//
//	transport := httpmock.NewReplayTransport(t, httpmock.ReplayConfig{
//		Cassette: "my_test_cassette",
//		Mode:     httpmock.ModeUnitTest,
//	})
//	defer transport.Close()
//
//	client := &http.Client{Transport: transport}
//	// Now use client for requests to any hosts - they'll be recorded/replayed
func NewReplayTransport(tester T, config ReplayConfig) *ReplayTransport {
	tester.Helper()

	// Create the recorder
	rec, err := createRecorder(tester, config)
	require.NoError(tester, err, "failed to create recorder")

	return &ReplayTransport{
		rec:          rec,
		t:            tester,
		cassetteName: config.Cassette,
	}
}

// RoundTrip implements http.RoundTripper. It either replays the request from
// the cassette or forwards it to the real endpoint and records the result.
func (rt *ReplayTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Let the recorder handle it (replay or forward+record)
	resp, err := rt.rec.RoundTrip(req)
	if err != nil {
		// Log the error details for debugging, similar to ReplayServer
		rt.t.Errorf("Replay transport error: %v", err)
		// Return the error for the HTTP client to handle
		return nil, fmt.Errorf("replay transport error: %w", err)
	}

	// Increment the used interactions counter
	rt.usedInteractions.Add(1)

	return resp, nil
}

// Assert verifies that all recorded interactions were used. This can be called
// separately from Close() to perform assertions without stopping the recorder.
func (rt *ReplayTransport) assert() error {
	// Get the cassette to check if all interactions were used
	cassette, err := cassette.Load(rt.cassetteName)
	if err != nil {
		return fmt.Errorf("failed to load cassette for verification: %w", err)
	}

	// Check if any interactions were not used
	usedCount := rt.usedInteractions.Load()
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

// Close stops the ReplayTransport and verifies that all recorded interactions
// were used. This should be called when the test is complete, typically in
// a defer statement.
func (rt *ReplayTransport) Close() {
	// Stop the recorder
	err := rt.rec.Stop()
	require.NoError(rt.t, err, "failed to stop recorder")

	// Verify that all recorded interactions were used
	err = rt.assert()
	require.NoError(rt.t, err, "not all recorded interactions were used")
}

// NewReplayClient creates a pre-configured http.Client that records and replays
// HTTP interactions. This is a convenience function that creates a ReplayTransport
// and wraps it in an http.Client.
//
// Example usage:
//
//	client, close := httpmock.NewReplayClient(t, httpmock.ReplayConfig{
//		Cassette: "my_test_cassette",
//	})
//	defer close()
//
//	// Use client for requests to any hosts - they'll be recorded/replayed
//	resp, err := client.Get("https://api.example.com/users")
func NewReplayClient(tester T, config ReplayConfig) (*http.Client, func()) {
	tester.Helper()

	transport := NewReplayTransport(tester, config)

	client := &http.Client{Transport: transport}
	return client, transport.Close
}
