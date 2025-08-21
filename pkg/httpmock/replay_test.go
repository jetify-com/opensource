package httpmock

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

func TestNewReplayServer(t *testing.T) {
	// Test valid configuration
	replayServer := NewReplayServer(t, ReplayConfig{
		Host:     "https://jsonplaceholder.typicode.com",
		Cassette: "testdata/server_valid",
	})
	require.NotNil(t, replayServer, "replay server should not be nil")
	assert.Equal(t, "https://jsonplaceholder.typicode.com", replayServer.realURL.String())

	// Test cleanup
	replayServer.Close()
}

func TestReplayServerHandler(t *testing.T) {
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

			// Create a server with a normal request and response
			server := NewReplayServer(mockTester, ReplayConfig{
				Host:     "https://httpbin.org",
				Cassette: "testdata/server_" + test.name,
			})
			defer server.Close()

			// Build the request using the declarative Request struct
			req, err := buildRequest(server.URL(), test.request)
			require.NoError(t, err)

			// Make the request
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			// Check response status
			assert.Equal(t, test.expectedStatus, resp.StatusCode)

			// Check response body
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), test.expectedBody)

			// Close the server to ensure the cassette is saved
			server.Close()
		})
	}
}

func TestFormatRequestBody(t *testing.T) {
	tests := []struct {
		name     string
		body     io.ReadCloser
		expected string
	}{
		{
			name:     "nil body",
			body:     nil,
			expected: "",
		},
		{
			name:     "empty body",
			body:     io.NopCloser(bytes.NewReader([]byte{})),
			expected: "",
		},
		{
			name:     "string body",
			body:     io.NopCloser(bytes.NewReader([]byte("test body"))),
			expected: "test body",
		},
		{
			name:     "json body",
			body:     io.NopCloser(bytes.NewReader([]byte(`{"key":"value"}`))),
			expected: `{"key":"value"}`,
		},
		{
			name:     "error reading body",
			body:     errorReader{},
			expected: "error reading body: simulated read error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &http.Request{
				Body: test.body,
			}
			result := formatRequestBody(req)
			assert.Equal(t, test.expected, result)

			// If the body was readable, verify it can still be read
			if test.body != nil && test.name != "error reading body" {
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.Equal(t, test.expected, string(body), "body should be readable after formatting")
			}
		})
	}
}

func TestReplayServerAdditionalIgnoredHeaders(t *testing.T) {
	tests := []struct {
		name                     string
		additionalIgnoredHeaders []string
		requestHeaders           map[string]string
		responseHeaders          map[string]string
		expectedInCassette       struct {
			requestHeaders  map[string]string
			responseHeaders map[string]string
		}
	}{
		{
			name: "additional_ignored_headers",
			additionalIgnoredHeaders: []string{
				"X-Request-ID",
				"X-Trace-ID",
				"X-Response-ID",
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer token123",  // Default ignored
				"X-Request-ID":  "req-123",          // Additional ignored
				"X-Trace-ID":    "trace-456",        // Additional ignored
				"Content-Type":  "application/json", // Should be preserved
			},
			responseHeaders: map[string]string{
				"Set-Cookie":    "session=abc123", // Default ignored
				"X-Response-ID": "resp-789",       // Additional ignored
				"Server":        "test-server",    // Should be preserved
			},
			expectedInCassette: struct {
				requestHeaders  map[string]string
				responseHeaders map[string]string
			}{
				requestHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				responseHeaders: map[string]string{
					"Server":       "test-server",
					"Content-Type": "text/plain; charset=utf-8", // Auto-added by Go HTTP server
					// Note: Date header is also auto-added but varies, so we'll handle it separately
				},
			},
		},
		{
			name:                     "no_additional_ignored_headers",
			additionalIgnoredHeaders: []string{}, // No additional headers
			requestHeaders: map[string]string{
				"Authorization": "Bearer token123",  // Default ignored
				"X-Custom":      "custom-value",     // Should be preserved
				"Content-Type":  "application/json", // Should be preserved
			},
			responseHeaders: map[string]string{
				"Set-Cookie": "session=abc123", // Default ignored
				"X-Custom":   "response-value", // Should be preserved
			},
			expectedInCassette: struct {
				requestHeaders  map[string]string
				responseHeaders map[string]string
			}{
				requestHeaders: map[string]string{
					"X-Custom":     "custom-value",
					"Content-Type": "application/json",
				},
				responseHeaders: map[string]string{
					"X-Custom":     "response-value",
					"Content-Type": "text/plain; charset=utf-8", // Auto-added by Go HTTP server
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test server that returns the configured response headers
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for header, value := range test.responseHeaders {
					w.Header().Set(header, value)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "ok"}`))
			}))
			defer testServer.Close()

			cassetteName := "testdata/server_" + test.name

			// Create replay server with additional ignored headers
			replayServer := NewReplayServer(t, ReplayConfig{
				Host:                     testServer.URL,
				Cassette:                 cassetteName,
				AdditionalIgnoredHeaders: test.additionalIgnoredHeaders,
			})

			// Make request with the configured headers
			req, err := buildRequest(replayServer.URL(), Request{
				Method:  http.MethodGet,
				Path:    "/test",
				Headers: test.requestHeaders,
			})
			require.NoError(t, err)

			// Make the request
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			resp.Body.Close()

			// Close server to ensure cassette is saved
			replayServer.Close()

			// Load and inspect the cassette
			cassette, err := cassette.Load(cassetteName)
			require.NoError(t, err)
			require.Len(t, cassette.Interactions, 1, "should have exactly one interaction")

			interaction := cassette.Interactions[0]

			// Verify request headers: cassette should contain exactly the expected headers
			for header, expectedValue := range test.expectedInCassette.requestHeaders {
				assert.Equal(t, expectedValue, interaction.Request.Headers.Get(header),
					"request header %q should be preserved in cassette", header)
			}
			// Verify no unexpected request headers are present
			for header := range interaction.Request.Headers {
				if _, expected := test.expectedInCassette.requestHeaders[header]; !expected {
					assert.Empty(t, interaction.Request.Headers.Get(header),
						"unexpected request header %q found in cassette", header)
				}
			}

			// Verify response headers: cassette should contain exactly the expected headers
			for header, expectedValue := range test.expectedInCassette.responseHeaders {
				assert.Equal(t, expectedValue, interaction.Response.Headers.Get(header),
					"response header %q should be preserved in cassette", header)
			}
			// Verify no unexpected response headers are present (except for auto-generated ones like Date)
			for header := range interaction.Response.Headers {
				if _, expected := test.expectedInCassette.responseHeaders[header]; !expected {
					// Allow certain headers that are automatically added by Go's HTTP server
					if header == "Date" {
						continue // Date header is auto-generated and varies
					}
					assert.Empty(t, interaction.Response.Headers.Get(header),
						"unexpected response header %q found in cassette", header)
				}
			}
		})
	}
}

func TestRemoveIgnored(t *testing.T) {
	tests := []struct {
		name        string
		interaction *cassette.Interaction
		wantHeaders http.Header
	}{
		{
			name: "removes all ignored headers",
			interaction: &cassette.Interaction{
				Request: cassette.Request{
					Headers: http.Header{
						"Authorization":       []string{"Bearer token"},
						"Content-Type":        []string{"application/json"},
						"User-Agent":          []string{"test-agent"},
						"X-Custom":            []string{"value"},
						"Accept":              []string{"*/*"},
						"Cookie":              []string{"session=123"},
						"X-Forwarded-For":     []string{"1.2.3.4"},
						"X-Api-Key":           []string{"secret"},
						"X-Csrf-Token":        []string{"token"},
						"X-Requested-With":    []string{"XMLHttpRequest"},
						"X-Real-Ip":           []string{"1.2.3.4"},
						"X-Client-Ip":         []string{"1.2.3.4"},
						"Proxy-Authorization": []string{"Basic auth"},
						"Www-Authenticate":    []string{"Basic"},
						"Proxy-Authenticate":  []string{"Basic"},
						"Set-Cookie":          []string{"session=456"},
					},
				},
				Response: cassette.Response{
					Headers: http.Header{
						"Set-Cookie":   []string{"session=789"},
						"Content-Type": []string{"application/json"},
					},
				},
			},
			wantHeaders: http.Header{
				"Content-Type": {"application/json"},
				"X-Custom":     {"value"},
			},
		},
		{
			name: "keeps non-ignored headers",
			interaction: &cassette.Interaction{
				Request: cassette.Request{
					Headers: http.Header{
						"Content-Type": []string{"application/json"},
						"X-Custom":     {"value"},
					},
				},
				Response: cassette.Response{
					Headers: http.Header{
						"Content-Type": []string{"application/json"},
						"X-Custom":     {"value"},
					},
				},
			},
			wantHeaders: http.Header{
				"Content-Type": {"application/json"},
				"X-Custom":     {"value"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := removeIgnoredHeaders(test.interaction, ignoredHeaders)
			assert.NoError(t, err)

			// Check request headers
			assert.Equal(t, test.wantHeaders, test.interaction.Request.Headers)

			// Check response headers
			for _, h := range ignoredHeaders {
				assert.Empty(t, test.interaction.Response.Headers[h], "ignored header %q should be removed from response", h)
			}
		})
	}
}

func TestReplayServerFailures(t *testing.T) {
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
				Method: http.MethodPost, // Cassette has GET /get, but we'll send POST
				Path:   "/get",
				Headers: map[string]string{
					"Host": "httpbin.org",
				},
			},
			expectedError: true,
			errorContains: "method mismatch",
			cassette:      "server_successful_get",
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
			errorContains: "URL mismatch",
			cassette:      "server_successful_get",
		},
		{
			name: "mismatched_body",
			request: Request{
				Method: http.MethodPost,
				Path:   "/post",
				Headers: map[string]string{
					"Host":         "httpbin.org",
					"Content-Type": "application/json",
				},
				Body: `{"unexpected": "data"}`, // Cassette has empty body, we'll send one
			},
			expectedError: true,
			errorContains: "body mismatch",
			cassette:      "server_post_with_body",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use a mockT to capture failures
			mockTester := &mockT{}

			// Create a server with a normal request and response
			replayServer := NewReplayServer(mockTester, ReplayConfig{
				Host:     "https://httpbin.org",
				Cassette: filepath.Join("testdata", test.cassette),
			})
			defer replayServer.Close()

			// Build the request using the declarative Request struct
			req, err := buildRequest(replayServer.URL(), test.request)
			require.NoError(t, err)

			// Make the request - this should trigger the test failure via requireCassetteRequest
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			// Check if the test failed as expected
			if test.expectedError {
				assert.True(t, mockTester.failed, "test should have failed")
				if test.errorContains != "" {
					assert.Contains(t, strings.Join(mockTester.errors, "\n"), test.errorContains, "expected error message not found")
				}
			} else {
				assert.False(t, mockTester.failed, "test should not have failed")
			}
		})
	}
}

func TestReplayServerInteractionCounts(t *testing.T) {
	tests := []struct {
		name          string
		requests      []Request
		expectedError bool
		errorContains string
		cassette      string
		checkCloseErr bool // indicates whether to check error from Close() or mockT
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
			errorContains: "expected 2 requests, received 0. Next expected: [GET https://httpbin.org/get]",
			cassette:      "server_multiple_interactions",
			checkCloseErr: true,
		},
		{
			name: "too_many_requests",
			requests: []Request{
				{
					Method: http.MethodGet,
					Path:   "/get",
					Headers: map[string]string{
						"Host": "httpbin.org",
					},
				},
				{
					Method: http.MethodGet,
					Path:   "/get",
					Headers: map[string]string{
						"Host": "httpbin.org",
					},
				},
			},
			expectedError: true,
			errorContains: "requested interaction not found",
			cassette:      "server_successful_get",
			checkCloseErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use a mockT to capture failures
			mockTester := &mockT{}

			// Create a server with a normal request and response
			replayServer := NewReplayServer(mockTester, ReplayConfig{
				Host:     "https://httpbin.org",
				Cassette: filepath.Join("testdata", test.cassette),
			})
			defer replayServer.Close()

			// Make all the requests
			for _, request := range test.requests {
				req, err := buildRequest(replayServer.URL(), request)
				require.NoError(t, err)
				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				defer func() { _ = resp.Body.Close() }()
			}

			// Close the server and check for errors
			replayServer.Close()

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
