package httpmock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReplayTransport(t *testing.T) {
	// Test valid configuration
	transport := NewReplayTransport(t, ReplayConfig{
		Cassette: "testdata/transport_valid",
	})
	require.NotNil(t, transport, "transport should not be nil")
	assert.Equal(t, "testdata/transport_valid", transport.cassetteName)

	// Test cleanup
	transport.Close()
}

func TestReplayTransportRoundTrip(t *testing.T) {
	tests := []struct {
		name           string
		request        Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful_get",
			request: Request{
				Method: http.MethodGet,
				Path:   "/get",
				Headers: map[string]string{
					"Host": "httpbin.org",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"url": "https://httpbin.org/get"`,
		},
		{
			name: "post_with_body",
			request: Request{
				Method: http.MethodPost,
				Path:   "/post",
				Headers: map[string]string{
					"Host":         "httpbin.org",
					"Content-Type": "application/json",
				},
				Body: `{"test": "data"}`,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"test": "data"`,
		},
		{
			name: "get_with_header",
			request: Request{
				Method: http.MethodGet,
				Path:   "/headers",
				Headers: map[string]string{
					"Host":          "httpbin.org",
					"X-Test-Header": "test-value",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"X-Test-Header": "test-value"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a mock T that will record failures
			mockTester := &mockT{}

			// Create a transport
			transport := NewReplayTransport(mockTester, ReplayConfig{
				Cassette: "testdata/transport_" + test.name,
			})
			defer transport.Close()

			// Create HTTP client with our transport
			client := &http.Client{Transport: transport}

			// Build the request - note we need to use the actual target URL,
			// not a proxy URL like with ReplayServer
			req, err := buildRequestForHost("https://httpbin.org", test.request)
			require.NoError(t, err)

			// Make the request
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			assert.Equal(t, test.expectedStatus, resp.StatusCode)

			// Check response body
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), test.expectedBody)

			// Close the transport to ensure the cassette is saved
			transport.Close()
		})
	}
}

func TestReplayTransportMultipleHosts(t *testing.T) {
	// Create a mock T that will record failures
	mockTester := &mockT{}

	// Create a transport that can handle multiple hosts
	transport := NewReplayTransport(mockTester, ReplayConfig{
		Cassette: "testdata/transport_multi_host",
	})
	defer transport.Close()

	// Create HTTP client with our transport
	client := &http.Client{Transport: transport}

	// Make requests to different hosts
	hosts := []string{
		"https://httpbin.org",
		"https://jsonplaceholder.typicode.com",
	}

	for i, host := range hosts {
		t.Run(host, func(t *testing.T) {
			req, err := buildRequestForHost(host, Request{
				Method: http.MethodGet,
				Path:   "/get",
			})
			if host == "https://jsonplaceholder.typicode.com" {
				req, err = buildRequestForHost(host, Request{
					Method: http.MethodGet,
					Path:   "/posts/1",
				})
			}
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			// Should get successful responses
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify we can read the body
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Greater(t, len(body), 0, "response should have body")

			t.Logf("Request %d to %s: status %d, body length %d", i+1, host, resp.StatusCode, len(body))
		})
	}

	// Close the transport
	transport.Close()
}

func TestReplayTransportFailures(t *testing.T) {
	tests := []struct {
		name          string
		request       Request
		expectedError bool
		errorContains string
		cassette      string
	}{
		{
			name: "mismatched_method",
			request: Request{
				Method: http.MethodPost, // Cassette has GET, but we'll send POST
				Path:   "/get",
				Headers: map[string]string{
					"Host": "httpbin.org",
				},
			},
			expectedError: true,
			errorContains: "requested interaction not found",
			cassette:      "transport_successful_get",
		},
		{
			name: "mismatched_path",
			request: Request{
				Method: http.MethodGet,
				Path:   "/wrong", // Cassette has /get, but we'll send /wrong
				Headers: map[string]string{
					"Host": "httpbin.org",
				},
			},
			expectedError: true,
			errorContains: "requested interaction not found",
			cassette:      "transport_successful_get",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use a mockT to capture failures
			mockTester := &mockT{}

			// Create a transport
			transport := NewReplayTransport(mockTester, ReplayConfig{
				Cassette: "testdata/" + test.cassette,
			})
			defer transport.Close()

			// Create HTTP client with our transport
			client := &http.Client{Transport: transport}

			// Build the request
			req, err := buildRequestForHost("https://httpbin.org", test.request)
			require.NoError(t, err)

			// Make the request - transport returns direct errors, not test failures
			resp, err := client.Do(req)

			if test.expectedError {
				assert.Error(t, err, "expected request to fail")
				if test.errorContains != "" {
					assert.Contains(t, err.Error(), test.errorContains, "expected error message not found")
				}
			} else {
				assert.NoError(t, err, "request should not fail")
				if resp != nil {
					defer func() { _ = resp.Body.Close() }()
				}
			}
		})
	}
}

func TestReplayTransportInteractionCounts(t *testing.T) {
	tests := []struct {
		name          string
		requests      []Request
		expectedError bool
		errorContains string
		cassette      string
	}{
		{
			name: "too_few_requests",
			requests: []Request{
				{
					Method: http.MethodGet,
					Path:   "/get",
					Headers: map[string]string{
						"Host": "httpbin.org",
					},
				},
			},
			expectedError: true,
			errorContains: "expected 2 requests, received 1. Next expected: [GET https://httpbin.org/status/200]",
			cassette:      "server_multiple_interactions",
		},
		{
			name: "correct_number_of_requests",
			requests: []Request{
				{
					Method: http.MethodGet,
					Path:   "/get",
					Headers: map[string]string{
						"Host": "httpbin.org",
					},
				},
			},
			expectedError: false,
			cassette:      "transport_successful_get",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use a mockT to capture failures
			mockTester := &mockT{}

			// Create a transport
			transport := NewReplayTransport(mockTester, ReplayConfig{
				Cassette: "testdata/" + test.cassette,
			})
			defer transport.Close()

			// Create HTTP client with our transport
			client := &http.Client{Transport: transport}

			// Make all the requests
			for _, request := range test.requests {
				req, err := buildRequestForHost("https://httpbin.org", request)
				require.NoError(t, err)
				resp, err := client.Do(req)
				require.NoError(t, err)
				defer func() { _ = resp.Body.Close() }()
			}

			// Close the transport and check for errors
			transport.Close()

			// Verify the test failed as expected
			if test.expectedError {
				assert.True(t, mockTester.failed, "test should have failed")
				assert.Contains(t, strings.Join(mockTester.errors, "\n"), test.errorContains, "expected error message not found")
			} else {
				assert.False(t, mockTester.failed, "test should not have failed")
			}
		})
	}
}

// buildRequestForHost creates an HTTP request for a specific host
// This is different from the buildRequest helper used in replay_test.go
// which builds requests for the proxy server
func buildRequestForHost(baseURL string, req Request) (*http.Request, error) {
	url := baseURL + req.Path

	var body io.Reader
	if req.Body != nil {
		switch b := req.Body.(type) {
		case string:
			body = bytes.NewReader([]byte(b))
		case []byte:
			body = bytes.NewReader(b)
		default:
			return nil, fmt.Errorf("unsupported body type: %T", req.Body)
		}
	}

	httpReq, err := http.NewRequest(req.Method, url, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	return httpReq, nil
}

func TestNewReplayClient(t *testing.T) {
	// Test valid configuration
	client, close := NewReplayClient(t, ReplayConfig{
		Cassette: "testdata/transport_client_test",
	})
	require.NotNil(t, client, "client should not be nil")
	require.NotNil(t, close, "close function should not be nil")

	// Test that we can use the client
	resp, err := client.Get("https://httpbin.org/get")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test cleanup
	close()
}

func TestReplayClientMultipleHosts(t *testing.T) {
	// Create a client that can handle multiple hosts
	client, close := NewReplayClient(t, ReplayConfig{
		Cassette: "testdata/transport_client_multi_host",
	})
	defer close()

	// Make requests to different hosts using the same client
	hosts := []string{
		"https://httpbin.org/get",
		"https://jsonplaceholder.typicode.com/posts/1",
	}

	for i, hostURL := range hosts {
		t.Run(fmt.Sprintf("request_%d", i+1), func(t *testing.T) {
			resp, err := client.Get(hostURL)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			// Should get successful responses
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify we can read the body
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Greater(t, len(body), 0, "response should have body")

			t.Logf("Request to %s: status %d, body length %d", hostURL, resp.StatusCode, len(body))
		})
	}

	// Close the client
	close()
}
