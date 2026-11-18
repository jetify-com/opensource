package httpmock_test

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/require"
	"go.jetify.com/pkg/httpmock"
)

// ExampleReplayTransport demonstrates using ReplayTransport to record and replay
// HTTP interactions across multiple hosts. This is ideal for testing SDKs that
// communicate with multiple different services.
func ExampleReplayTransport() {
	mockT := &t{} // In a real test, this would be the actual testing.T instance

	// Create a ReplayTransport that can record/replay requests to any host
	transport := httpmock.NewReplayTransport(mockT, httpmock.ReplayConfig{
		Cassette: "testdata/transport_multi_host_example",
	})
	defer transport.Close()

	// Create an HTTP client that uses our transport
	client := &http.Client{Transport: transport}

	// Now the client can make requests to different hosts, and they'll all be
	// recorded/replayed from the same cassette

	// Request to first service
	resp1, err := client.Get("https://httpbin.org/get")
	require.NoError(mockT, err)
	require.NoError(mockT, resp1.Body.Close())

	// Request to second service
	resp2, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	require.NoError(mockT, err)
	require.NoError(mockT, resp2.Body.Close())

	fmt.Printf("First service status: %d\n", resp1.StatusCode)
	fmt.Printf("Second service status: %d\n", resp2.StatusCode)

	// Output:
	// First service status: 200
	// Second service status: 200
}

// ExampleNewReplayClient demonstrates the convenience function for creating
// a pre-configured HTTP client for recording and replaying requests.
func ExampleNewReplayClient() {
	mockT := &t{} // In a real test, this would be the actual testing.T instance

	// Create a pre-configured HTTP client - one line!
	client, close := httpmock.NewReplayClient(mockT, httpmock.ReplayConfig{
		Cassette: "testdata/transport_client_example",
	})
	defer close()

	// Use the client directly - no transport configuration needed
	resp, err := client.Get("https://httpbin.org/get")
	require.NoError(mockT, err)
	defer func() { _ = resp.Body.Close() }()

	fmt.Printf("Response status: %d\n", resp.StatusCode)

	// Output:
	// Response status: 200
}
