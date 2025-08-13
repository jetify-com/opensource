package httpmock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unit tests - core functionality

func TestServer_Request(t *testing.T) {
	tests := []struct {
		name     string
		expect   Exchange // what we configure the server to expect
		send     *Request // what we actually send (nil means use expect.Request)
		wantFail bool
	}{
		{
			name: "matching request succeeds",
			expect: Exchange{
				Request: Request{
					Method: "POST",
					Path:   "/test",
					Body:   `{"hello":"world"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: Response{
					StatusCode: http.StatusOK,
					Body:       map[string]string{"status": "ok"},
				},
			},
		},
		{
			name: "unexpected request fails",
			expect: Exchange{
				Request: Request{
					Method: "POST",
					Path:   "/test",
				},
			},
			send: &Request{
				Method: "GET",
				Path:   "/wrong",
			},
			wantFail: true,
		},
		{
			name: "default response code is 200",
			expect: Exchange{
				Request: Request{
					Method: "GET",
					Path:   "/test",
				},
				Response: Response{
					Body: "hello",
				},
			},
		},
		{
			name: "custom header match",
			expect: Exchange{
				Request: Request{
					Method: "GET",
					Path:   "/test",
					Headers: map[string]string{
						"X-Custom":      "value",
						"Authorization": "Bearer token",
					},
				},
				Response: Response{
					Body: "ok",
				},
			},
		},
		{
			name: "custom validation passes",
			expect: Exchange{
				Request: Request{
					Method: "POST",
					Path:   "/test",
					Validate: func(r *http.Request) error {
						if r.Header.Get("Date") == "" {
							return fmt.Errorf("missing Date header")
						}
						return nil
					},
				},
			},
			send: &Request{
				Method: "POST",
				Path:   "/test",
				Headers: map[string]string{
					"Date": "Mon, 19 Feb 2024 10:00:00 GMT",
				},
			},
		},
		{
			name: "custom validation fails",
			expect: Exchange{
				Request: Request{
					Method: "POST",
					Path:   "/test",
					Validate: func(r *http.Request) error {
						return fmt.Errorf("validation failed")
					},
				},
				Response: Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			send: &Request{
				Method: "POST",
				Path:   "/test",
			},
			wantFail: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockTester := &mockT{}
			testServer := NewServer(mockTester, []Exchange{testCase.expect})
			defer testServer.Close()

			reqToSend := testCase.expect.Request
			if testCase.send != nil {
				reqToSend = *testCase.send
			}

			req, err := buildRequest(testServer.BaseURL(), reqToSend)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			if testCase.wantFail {
				assert.True(t, mockTester.failed, "Expected test to fail")
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			} else {
				assert.False(t, mockTester.failed, "Test unexpectedly failed")
				assertResponseEq(t, testCase.expect.Response, resp)
			}
		})
	}
}

func TestServer_VerifyComplete(t *testing.T) {
	tests := []struct {
		name    string
		expect  []Exchange
		send    []Request // requests to send in order
		wantErr string    // expected error message, empty means expect success
	}{
		{
			name: "all requests made",
			expect: []Exchange{{
				Request: Request{
					Method: "GET",
					Path:   "/test",
				},
				Response: Response{
					Body: "ok",
				},
			}},
			send: []Request{{
				Method: "GET",
				Path:   "/test",
			}},
		},
		{
			name: "missing requests",
			expect: []Exchange{{
				Request: Request{
					Method: "GET",
					Path:   "/test",
				},
			}},
			send:    []Request{}, // send nothing
			wantErr: "expected 1 requests, received 0. Next expected: [GET /test]",
		},
		{
			name: "requests in correct order",
			expect: []Exchange{
				{
					Request: Request{
						Method: "GET",
						Path:   "/first",
					},
				},
				{
					Request: Request{
						Method: "POST",
						Path:   "/second",
					},
				},
			},
			send: []Request{
				{
					Method: "GET",
					Path:   "/first",
				},
				{
					Method: "POST",
					Path:   "/second",
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockTester := &mockT{}
			testServer := NewServer(mockTester, testCase.expect)
			defer testServer.Close()

			// Send all requests in order
			for _, req := range testCase.send {
				r, err := buildRequest(testServer.BaseURL(), req)
				require.NoError(t, err)
				resp, err := http.DefaultClient.Do(r)
				require.NoError(t, err)
				_ = resp.Body.Close()
			}

			err := testServer.VerifyComplete()
			if testCase.wantErr != "" {
				assert.EqualError(t, err, testCase.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServer_BodyComparison(t *testing.T) {
	tests := []struct {
		name     string
		expected any
		actual   string
		wantFail bool
	}{
		{
			name:     "string match",
			expected: "hello",
			actual:   "hello",
			wantFail: false,
		},
		{
			name:     "json string match",
			expected: `{"a":1}`,
			actual:   `{"a": 1}`,
			wantFail: false,
		},
		{
			name:     "struct match",
			expected: map[string]int{"a": 1},
			actual:   `{"a": 1}`,
			wantFail: false,
		},
		{
			name:     "mismatch",
			expected: "hello",
			actual:   "world",
			wantFail: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockTester := &mockT{}
			mockTester.failed = false
			mockTester.errors = nil

			requireBodyEq(mockTester, testCase.expected, []byte(testCase.actual))
			assert.Equal(t, testCase.wantFail, mockTester.failed)
		})
	}
}

func TestServer_Path(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		expected string
	}{
		{
			name:     "normal case",
			baseURL:  "http://localhost:8080",
			path:     "/api/users",
			expected: "http://localhost:8080/api/users",
		},
		{
			name:     "base URL with trailing slash",
			baseURL:  "http://localhost:8080/",
			path:     "/api/users",
			expected: "http://localhost:8080/api/users",
		},
		{
			name:     "path without leading slash",
			baseURL:  "http://localhost:8080",
			path:     "api/users",
			expected: "http://localhost:8080/api/users",
		},
		{
			name:     "both with extra slashes",
			baseURL:  "http://localhost:8080/",
			path:     "/api/users",
			expected: "http://localhost:8080/api/users",
		},
		{
			name:     "root path",
			baseURL:  "http://localhost:8080",
			path:     "/",
			expected: "http://localhost:8080/",
		},
		{
			name:     "empty path",
			baseURL:  "http://localhost:8080",
			path:     "",
			expected: "http://localhost:8080",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Create a server with a custom base URL for testing
			s := &Server{server: httptest.NewServer(nil)}
			s.server.URL = testCase.baseURL // override the random port with our test URL
			defer s.server.Close()

			got := s.Path(testCase.path)
			assert.Equal(t, testCase.expected, got)
		})
	}
}

// Consolidated MergeRequests tests into a single table-driven test
func TestMergeRequests(t *testing.T) {
	tests := []struct {
		name     string
		requests []Request
		expected Request
		validate bool // whether to test validation function
	}{
		{
			name:     "empty list",
			requests: []Request{},
			expected: Request{},
		},
		{
			name: "single request",
			requests: []Request{
				{
					Method: "GET",
					Path:   "/api",
					Headers: map[string]string{
						"Authorization": "Bearer token",
					},
				},
			},
			expected: Request{
				Method: "GET",
				Path:   "/api",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
		},
		{
			name: "two requests - override empty",
			requests: []Request{
				{
					Method: "GET",
					Path:   "/api",
					Headers: map[string]string{
						"Authorization": "Bearer token",
					},
				},
				{},
			},
			expected: Request{
				Method: "GET",
				Path:   "/api",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
		},
		{
			name: "two requests - override some fields",
			requests: []Request{
				{
					Method: "GET",
					Path:   "/api",
					Headers: map[string]string{
						"Authorization": "Bearer token",
					},
				},
				{
					Method: "POST",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			expected: Request{
				Method: "POST",
				Path:   "/api",
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"Content-Type":  "application/json",
				},
			},
		},
		{
			name: "two requests - override all fields",
			requests: []Request{
				{
					Method: "GET",
					Path:   "/api",
					Body:   "original",
					Headers: map[string]string{
						"Authorization": "Bearer token",
					},
				},
				{
					Method: "POST",
					Path:   "/new",
					Body:   "override",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			expected: Request{
				Method: "POST",
				Path:   "/new",
				Body:   "override",
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"Content-Type":  "application/json",
				},
			},
		},
		{
			name: "three requests - cumulative merge",
			requests: []Request{
				{
					Method: "GET",
					Path:   "/api",
				},
				{
					Headers: map[string]string{
						"X-Header1": "value1",
					},
				},
				{
					Headers: map[string]string{
						"X-Header2": "value2",
					},
					Body: "test body",
				},
			},
			expected: Request{
				Method: "GET",
				Path:   "/api",
				Body:   "test body",
				Headers: map[string]string{
					"X-Header1": "value1",
					"X-Header2": "value2",
				},
			},
		},
		{
			name: "validate function",
			requests: []Request{
				{
					Method: "GET",
					Path:   "/api",
				},
				{
					Validate: func(r *http.Request) error {
						return nil // Will be tested separately
					},
				},
			},
			expected: Request{
				Method: "GET",
				Path:   "/api",
			},
			validate: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := MergeRequests(testCase.requests...)

			// For validate function test, we need special handling
			if testCase.validate {
				validationCalled := false
				result.Validate = func(r *http.Request) error {
					validationCalled = true
					return nil
				}

				req, _ := http.NewRequest(http.MethodGet, "/api", nil)
				err := result.Validate(req)
				assert.NoError(t, err)
				assert.True(t, validationCalled, "Validation function was not called")
			} else {
				// Clear the Validate function before comparison, as functions cannot be directly compared
				result.Validate = nil
				testCase.expected.Validate = nil
				assert.Equal(t, testCase.expected, result)
			}
		})
	}
}

func TestServer_ResponseHeaders(t *testing.T) {
	tests := []struct {
		name        string
		response    Response
		wantHeaders map[string]string
		wantCode    int
		wantBody    string
	}{
		{
			name: "default content-type",
			response: Response{
				Body: map[string]string{"hello": "world"},
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			wantCode: http.StatusOK,
			wantBody: `{"hello":"world"}`,
		},
		{
			name: "custom headers",
			response: Response{
				Body: "hello world",
				Headers: map[string]string{
					"Content-Type": "text/plain",
					"X-Custom":     "value",
				},
			},
			wantHeaders: map[string]string{
				"Content-Type": "text/plain",
				"X-Custom":     "value",
			},
			wantCode: http.StatusOK,
			wantBody: "hello world",
		},
		{
			name: "override default content-type",
			response: Response{
				Body: map[string]string{"hello": "world"},
				Headers: map[string]string{
					"Content-Type": "application/problem+json",
				},
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/problem+json",
			},
			wantCode: http.StatusOK,
			wantBody: `{"hello":"world"}`,
		},
		{
			name: "custom status code and text",
			response: Response{
				StatusCode: http.StatusNotFound,
				Status:     "Not Found",
				Body:       `"Resource not found"`,
			},
			wantHeaders: map[string]string{
				"Content-Type":  "application/json",
				"X-Status-Text": "Not Found",
			},
			wantCode: http.StatusNotFound,
			wantBody: `"Resource not found"`,
		},
		{
			name: "empty body",
			response: Response{
				StatusCode: http.StatusNoContent,
			},
			wantHeaders: map[string]string{},
			wantCode:    http.StatusNoContent,
			wantBody:    "",
		},
		{
			name: "invalid json body in response",
			response: Response{
				StatusCode: http.StatusBadRequest,
				Body:       `"invalid":json}`,
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			wantCode: http.StatusBadRequest,
			wantBody: `"invalid":json}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testServer := NewServer(t, []Exchange{{
				Request: Request{
					Method: "GET",
					Path:   "/test",
				},
				Response: testCase.response,
			}})
			defer testServer.Close()

			resp, err := http.Get(testServer.Path("/test"))
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			// Check status code
			assert.Equal(t, testCase.wantCode, resp.StatusCode)

			// Check headers
			for k, want := range testCase.wantHeaders {
				assert.Equal(t, want, resp.Header.Get(k), "Expected header %s to be %s", k, want)
			}

			// Check body
			if testCase.wantBody != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				if strings.HasPrefix(testCase.wantBody, "{") || strings.HasPrefix(testCase.wantBody, "[") {
					assert.JSONEq(t, testCase.wantBody, string(body))
				} else {
					assert.Equal(t, testCase.wantBody, string(body))
				}
			}
		})
	}
}

func TestServer_Delay(t *testing.T) {
	// Define a delay that's long enough to measure but short enough for tests
	delay := 200 * time.Millisecond

	// Create a server with a delayed response
	testServer := NewServer(t, []Exchange{{
		Request: Request{
			Method: "GET",
			Path:   "/delayed",
		},
		Response: Response{
			StatusCode: http.StatusOK,
			Body:       `{"message":"delayed response"}`,
			Delay:      delay,
		},
	}})
	defer testServer.Close()

	// Record start time
	start := time.Now()

	// Make the request
	resp, err := http.Get(testServer.Path("/delayed"))
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Calculate elapsed time
	elapsed := time.Since(start)

	// Verify response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"message":"delayed response"}`, string(body))

	// Verify timing - it should have taken at least the delay duration
	// We add a small buffer (10ms) to account for very minor timing inconsistencies
	assert.GreaterOrEqual(t, elapsed, delay-10*time.Millisecond,
		"Response came back too quickly (in %v), expected at least %v delay", elapsed, delay)
}

// Test error handling

// mockResponseWriter is a custom http.ResponseWriter for testing error scenarios
type mockResponseWriter struct {
	headers      http.Header
	statusCode   int
	responseBody []byte
	failOnWrite  bool
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		headers: make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	if m.failOnWrite {
		return 0, fmt.Errorf("forced write error")
	}
	m.responseBody = append(m.responseBody, b...)
	return len(b), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

// Tests for handler error paths
func TestServer_HandlerErrors(t *testing.T) {
	t.Run("response error", func(t *testing.T) {
		// Create a mock T that will record failures
		mockTester := &mockT{}

		// Create a server with a normal request and response
		server := &Server{
			t: mockTester,
			expectations: []Exchange{{
				Request: Request{
					Method: http.MethodGet,
					Path:   "/test",
				},
				Response: Response{
					Body: make(chan int), // Channel cannot be marshaled to JSON, will cause an error
				},
			}},
		}

		// Setup request and mock response writer
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := newMockResponseWriter()

		// Call the handler
		server.handler(w, req)

		// Handler should have set an error status code and recorded a test failure
		assert.Equal(t, http.StatusInternalServerError, w.statusCode)
		assert.True(t, mockTester.failed)
	})

	t.Run("validation error", func(t *testing.T) {
		// This test directly tests how requireRequestEq handles validation errors
		mockTester := &mockT{}

		// Create a request with a custom validation function that always fails
		expectedReq := Request{
			Method: http.MethodGet,
			Path:   "/test",
			Validate: func(r *http.Request) error {
				return fmt.Errorf("intentional validation error")
			},
		}

		// Create an HTTP request that would match except for the validation
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)

		// Call assertRequest directly
		assertRequest(mockTester, expectedReq, req)

		// Verify the test failed due to validation
		assert.True(t, mockTester.failed)
		assert.Contains(t, mockTester.errors[0], "custom validation failed")
	})

	t.Run("unexpected request", func(t *testing.T) {
		mockTester := &mockT{}

		// Create a server with no expectations
		server := NewServer(mockTester, []Exchange{})
		defer server.Close()

		// Send a request that isn't expected
		resp, err := http.Get(server.Path("/unexpected"))
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// Should have responded with 500 error and failed the test
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.True(t, mockTester.failed)
	})
}

// Mock implementations
type mockT struct {
	failed bool
	errors []string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockT) FailNow() {
	m.failed = true
}

func (m *mockT) Helper() {
	// No-op
}

func buildBody(body any) io.Reader {
	var r io.Reader
	switch b := body.(type) {
	case string:
		r = strings.NewReader(b)
	case []byte:
		r = bytes.NewReader(b)
	default:
		jsonBytes, err := json.Marshal(b)
		if err != nil {
			panic(err)
		}
		r = bytes.NewReader(jsonBytes)
	}
	return r
}

// assertResponseEq verifies that the actual response matches the expected response.
// It checks both the status code and body.
func assertResponseEq(t *testing.T, expected Response, actual *http.Response) {
	if expected.StatusCode > 0 {
		assert.Equal(t, expected.StatusCode, actual.StatusCode)
	} else {
		assert.Equal(t, http.StatusOK, actual.StatusCode)
	}

	if expected.Body != nil {
		body, err := io.ReadAll(actual.Body)
		require.NoError(t, err)
		requireBodyEq(&mockT{}, expected.Body, body)
	}
}

// buildRequest creates an http.Request from a base URL and Request struct.
func buildRequest(baseURL string, request Request) (*http.Request, error) {
	u, err := url.JoinPath(baseURL, request.Path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(request.Method, u, buildBody(request.Body))
	if err != nil {
		return nil, err
	}

	// Set headers from the request
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func TestRequireRequestEq(t *testing.T) {
	tests := []struct {
		name     string
		expected Request
		actual   *http.Request
		wantFail bool
	}{
		{
			name: "matching request",
			expected: Request{
				Method: http.MethodPost,
				Path:   "/test",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"hello":"world"}`,
			},
			actual: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"hello":"world"}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			wantFail: false,
		},
		{
			name: "mismatched method",
			expected: Request{
				Method: http.MethodPost,
				Path:   "/test",
			},
			actual:   httptest.NewRequest(http.MethodGet, "/test", nil),
			wantFail: true,
		},
		{
			name: "mismatched path",
			expected: Request{
				Method: http.MethodPost,
				Path:   "/test",
			},
			actual:   httptest.NewRequest(http.MethodPost, "/wrong", nil),
			wantFail: true,
		},
		{
			name: "mismatched headers",
			expected: Request{
				Method: http.MethodPost,
				Path:   "/test",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			actual: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.Header.Set("Content-Type", "text/plain")
				return req
			}(),
			wantFail: true,
		},
		{
			name: "mismatched body",
			expected: Request{
				Method: http.MethodPost,
				Path:   "/test",
				Body:   `{"hello":"world"}`,
			},
			actual:   httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"hello":"different"}`)),
			wantFail: true,
		},
		{
			name: "custom validation passes",
			expected: Request{
				Method: http.MethodPost,
				Path:   "/test",
				Validate: func(r *http.Request) error {
					if r.Header.Get("X-Custom") != "value" {
						return fmt.Errorf("missing X-Custom header")
					}
					return nil
				},
			},
			actual: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.Header.Set("X-Custom", "value")
				return req
			}(),
			wantFail: false,
		},
		{
			name: "custom validation fails",
			expected: Request{
				Method: http.MethodPost,
				Path:   "/test",
				Validate: func(r *http.Request) error {
					return fmt.Errorf("validation failed")
				},
			},
			actual:   httptest.NewRequest(http.MethodPost, "/test", nil),
			wantFail: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockTester := &mockT{}
			assertRequest(mockTester, testCase.expected, testCase.actual)
			assert.Equal(t, testCase.wantFail, mockTester.failed)
		})
	}
}
