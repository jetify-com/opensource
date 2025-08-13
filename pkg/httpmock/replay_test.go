package httpmock

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

func TestNewReplayServer(t *testing.T) {
	tests := []struct {
		name        string
		config      ReplayConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "valid configuration",
			config: ReplayConfig{
				Host:     "https://jsonplaceholder.typicode.com",
				Cassette: filepath.Join("testdata", "valid"),
			},
			wantErr: false,
		},
		{
			name: "invalid host URL",
			config: ReplayConfig{
				Host:     "://invalid-url",
				Cassette: filepath.Join("testdata", "invalid"),
			},
			wantErr:     true,
			errContains: "invalid Host",
		},
		{
			name: "empty host",
			config: ReplayConfig{
				Host:     "",
				Cassette: filepath.Join("testdata", "empty"),
			},
			wantErr:     true,
			errContains: "invalid Host",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			replayServer, err := NewReplayServer(t, test.config)
			if test.wantErr {
				assert.Error(t, err)
				if test.errContains != "" {
					assert.Contains(t, err.Error(), test.errContains)
				}
				assert.Nil(t, replayServer)
				return
			}
			require.NotNil(t, replayServer, "replay server should not be nil")
			assert.Equal(t, test.config.Host, replayServer.realURL.String())

			// Test cleanup
			err = replayServer.Close()
			assert.NoError(t, err)
		})
	}
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
			server, err := NewReplayServer(mockTester, ReplayConfig{
				Host:     "https://httpbin.org",
				Cassette: "testdata/" + test.name,
			})
			require.NoError(t, err)
			defer func() { _ = server.Close() }()

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
			err = server.Close()
			require.NoError(t, err)
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
			err := removeIgnored(test.interaction)
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
			cassette:      "successful_get",
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
			cassette:      "successful_get",
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
			cassette:      "post_with_body",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use a mockT to capture failures
			mockTester := &mockT{}

			// Create a server with a normal request and response
			replayServer, err := NewReplayServer(mockTester, ReplayConfig{
				Host:     "https://httpbin.org",
				Cassette: filepath.Join("testdata", test.cassette),
			})
			require.NoError(t, err)
			defer func() { _ = replayServer.Close() }()

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
			cassette:      "multiple_interactions",
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
			cassette:      "successful_get",
			checkCloseErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use a mockT to capture failures
			mockTester := &mockT{}

			// Create a server with a normal request and response
			replayServer, err := NewReplayServer(mockTester, ReplayConfig{
				Host:     "https://httpbin.org",
				Cassette: filepath.Join("testdata", test.cassette),
			})
			require.NoError(t, err)
			defer func() { _ = replayServer.Close() }()

			// Make all the requests
			for _, request := range test.requests {
				req, err := buildRequest(replayServer.URL(), request)
				require.NoError(t, err)
				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				defer func() { _ = resp.Body.Close() }()
			}

			// Close the server and check for errors
			err = replayServer.Close()

			// Verify the test failed as expected
			if test.expectedError {
				if test.checkCloseErr {
					assert.Contains(t, err.Error(), test.errorContains, "expected error message not found")
				} else {
					assert.True(t, mockTester.failed, "test should have failed")
					assert.Contains(t, strings.Join(mockTester.errors, "\n"), test.errorContains, "expected error message not found")
				}
			} else {
				assert.False(t, mockTester.failed, "test should not have failed")
				assert.NoError(t, err)
			}
		})
	}
}
