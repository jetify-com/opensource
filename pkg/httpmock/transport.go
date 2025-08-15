package httpmock

import (
	"fmt"
	"net/http"
	"sync/atomic"

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
//	transport, err := httpmock.NewReplayTransport(t, httpmock.ReplayConfig{
//		Cassette: "my_test_cassette",
//		Mode:     httpmock.ModeUnitTest,
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer transport.Close()
//
//	client := &http.Client{Transport: transport}
//	// Now use client for requests to any hosts - they'll be recorded/replayed
func NewReplayTransport(tester T, config ReplayConfig) (*ReplayTransport, error) {
	tester.Helper()

	// Create the recorder
	rec, err := createRecorder(tester, config)
	if err != nil {
		return nil, err
	}

	return &ReplayTransport{
		rec:          rec,
		t:            tester,
		cassetteName: config.Cassette,
	}, nil
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

// Close stops the ReplayTransport and verifies that all recorded interactions
// were used. This should be called when the test is complete, typically in
// a defer statement.
func (rt *ReplayTransport) Close() error {
	// Stop the recorder
	if err := rt.rec.Stop(); err != nil {
		return err
	}

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

// NewReplayClient creates a pre-configured http.Client that records and replays
// HTTP interactions. This is a convenience function that creates a ReplayTransport
// and wraps it in an http.Client.
//
// Example usage:
//
//	client, close, err := httpmock.NewReplayClient(t, httpmock.ReplayConfig{
//		Cassette: "my_test_cassette",
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer close()
//
//	// Use client for requests to any hosts - they'll be recorded/replayed
//	resp, err := client.Get("https://api.example.com/users")
func NewReplayClient(tester T, config ReplayConfig) (*http.Client, func() error, error) {
	tester.Helper()

	transport, err := NewReplayTransport(tester, config)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{Transport: transport}
	return client, transport.Close, nil
}
