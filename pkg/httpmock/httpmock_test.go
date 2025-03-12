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

// 1. Examples
func ExampleServer_basic() {
	testServer := NewServer(&mockT{}, []Exchange{{
		Request: Request{
			Method: "GET",
			Path:   "/hello",
		},
		Response: Response{
			Body: "world",
		},
	}})
	defer testServer.Close()

	resp, _ := http.Get(testServer.Path("/hello"))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: world
}

func ExampleServer_jsonRequest() {
	testServer := NewServer(&mockT{}, []Exchange{{
		Request: Request{
			Method: "POST",
			Path:   "/api/users",
			Body:   `{"name":"Alice"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		Response: Response{
			StatusCode: http.StatusCreated,
			Body: map[string]interface{}{
				"id":   1,
				"name": "Alice",
			},
		},
	}})
	defer testServer.Close()

	resp, _ := http.Post(testServer.Path("/api/users"),
		"application/json",
		strings.NewReader(`{"name":"Alice"}`))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	// Output:
	// 201
	// {"id":1,"name":"Alice"}
}

func ExampleServer_sequence() {
	testServer := NewServer(&mockT{}, []Exchange{
		{
			Request: Request{
				Method: "POST",
				Path:   "/login",
				Body:   `{"username":"alice","password":"secret"}`,
			},
			Response: Response{
				Body: map[string]string{"token": "abc123"},
			},
		},
		{
			Request: Request{
				Method: "GET",
				Path:   "/profile",
				Headers: map[string]string{
					"Authorization": "Bearer abc123",
				},
			},
			Response: Response{
				Body: map[string]string{"name": "Alice"},
			},
		},
	})
	defer testServer.Close()

	// Login
	resp, _ := http.Post(testServer.Path("/login"),
		"application/json",
		strings.NewReader(`{"username":"alice","password":"secret"}`))
	var loginResp struct{ Token string }
	err := json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	resp.Body.Close()

	// Get profile using token
	req, _ := http.NewRequest(http.MethodGet, testServer.Path("/profile"), nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp, _ = http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: {"name":"Alice"}
}

func ExampleServer_validation() {
	testServer := NewServer(&mockT{}, []Exchange{{
		Request: Request{
			Method: "POST",
			Path:   "/upload",
			Validate: func(r *http.Request) error {
				if r.Header.Get("Content-Length") == "0" {
					return fmt.Errorf("empty request body")
				}
				return nil
			},
		},
		Response: Response{
			StatusCode: http.StatusOK,
			Body:       "uploaded",
		},
	}})
	defer testServer.Close()

	// Send non-empty request
	resp, _ := http.Post(testServer.Path("/upload"),
		"text/plain",
		strings.NewReader("some data"))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: uploaded
}

func ExampleServer_delay() {
	// Create a server with a delayed response
	testServer := NewServer(&mockT{}, []Exchange{{
		Request: Request{
			Method: "GET",
			Path:   "/api/slow",
		},
		Response: Response{
			StatusCode: http.StatusOK,
			Body:       `{"status":"success", "data":"worth the wait"}`,
			Delay:      1 * time.Second, // Response will be delayed by 1 second
		},
	}})
	defer testServer.Close()

	// In a real test, you would use this URL with your HTTP client
	// fmt.Println("API endpoint:", testServer.Path("/api/slow"))

	// Instead of printing the URL (which has a dynamic port), print something static
	fmt.Println("Configured delay:", 1*time.Second)

	// Make a request (not executed in this example)
	// resp, err := http.Get(testServer.Path("/api/slow"))
	// The response would be delayed by 1 second

	// Output:
	// Configured delay: 1s
}

// 2. Unit tests (core functionality first)
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
			mockTester := newMockT(t)
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
			defer resp.Body.Close()

			if testCase.wantFail {
				assert.True(t, mockTester.failed)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			} else {
				assert.False(t, mockTester.failed)
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
			mockTester := newMockT(t)
			testServer := NewServer(mockTester, testCase.expect)
			defer testServer.Close()

			// Send all requests in order
			for _, req := range testCase.send {
				r, err := buildRequest(testServer.BaseURL(), req)
				require.NoError(t, err)
				resp, err := http.DefaultClient.Do(r)
				require.NoError(t, err)
				resp.Body.Close()
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
			mockTester := newMockT(t)
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

// 4. Mock implementations
type mockT struct {
	t      *testing.T // real testing.T for test assertions
	failed bool
	errors []string
}

func newMockT(t *testing.T) *mockT {
	return &mockT{t: t}
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockT) FailNow() {
	m.failed = true
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
		requireBodyEq(&mockT{t: t}, expected.Body, body)
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

func TestMergeRequests_TwoRequests(t *testing.T) {
	tests := []struct {
		name     string
		base     Request
		override Request
		want     Request
	}{
		{
			name: "override empty request",
			base: Request{
				Method: "GET",
				Path:   "/api",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
			override: Request{},
			want: Request{
				Method: "GET",
				Path:   "/api",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
		},
		{
			name: "override some fields",
			base: Request{
				Method: "GET",
				Path:   "/api",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
			override: Request{
				Method: "POST",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			want: Request{
				Method: "POST",
				Path:   "/api",
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"Content-Type":  "application/json",
				},
			},
		},
		{
			name: "override all fields",
			base: Request{
				Method: "GET",
				Path:   "/api",
				Body:   "original",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
			override: Request{
				Method: "POST",
				Path:   "/new",
				Body:   "override",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			want: Request{
				Method: "POST",
				Path:   "/new",
				Body:   "override",
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"Content-Type":  "application/json",
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := MergeRequests(testCase.base, testCase.override)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestServer_ResponseHeaders(t *testing.T) {
	tests := []struct {
		name        string
		response    Response
		wantHeaders map[string]string
	}{
		{
			name: "default content-type",
			response: Response{
				Body: map[string]string{"hello": "world"},
			},
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
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
			defer resp.Body.Close()

			for k, want := range testCase.wantHeaders {
				assert.Equal(t, want, resp.Header.Get(k))
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
	defer resp.Body.Close()

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

func TestMergeRequests_EmptyList(t *testing.T) {
	result := MergeRequests()
	assert.Equal(t, Request{}, result)
}

func TestMergeRequests_ValidateFunction(t *testing.T) {
	validationCalled := false

	base := Request{
		Method: "GET",
		Path:   "/api",
	}

	override := Request{
		Validate: func(r *http.Request) error {
			validationCalled = true
			return nil
		},
	}

	result := MergeRequests(base, override)

	// Create a test request to pass to the validate function
	req, _ := http.NewRequest("GET", "/api", nil)

	// Call the validate function from the merged request
	err := result.Validate(req)

	// Assert that our custom validate function was called
	assert.NoError(t, err)
	assert.True(t, validationCalled)
}

func TestMergeRequests_ThreeRequests(t *testing.T) {
	req1 := Request{
		Method: "GET",
		Path:   "/api",
	}

	req2 := Request{
		Headers: map[string]string{
			"X-Header1": "value1",
		},
	}

	req3 := Request{
		Headers: map[string]string{
			"X-Header2": "value2",
		},
		Body: "test body",
	}

	result := MergeRequests(req1, req2, req3)

	expected := Request{
		Method: "GET",
		Path:   "/api",
		Body:   "test body",
		Headers: map[string]string{
			"X-Header1": "value1",
			"X-Header2": "value2",
		},
	}

	assert.Equal(t, expected, result)
}

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
func TestServer_HandlerResponseError(t *testing.T) {
	// Create a mock T that will record failures
	mockTester := newMockT(t)

	// Create a server with a normal request and response
	server := &Server{
		t: mockTester,
		expectations: []Exchange{{
			Request: Request{
				Method: "GET",
				Path:   "/test",
			},
			Response: Response{
				Body: make(chan int), // Channel cannot be marshaled to JSON, will cause an error
			},
		}},
	}

	// Setup request and mock response writer
	req := httptest.NewRequest("GET", "/test", nil)
	w := newMockResponseWriter()

	// Call the handler
	server.handler(w, req)

	// Handler should have set an error status code and recorded a test failure
	assert.Equal(t, http.StatusInternalServerError, w.statusCode)
	assert.True(t, mockTester.failed)
}

func TestServer_HandlerValidationError(t *testing.T) {
	// This test directly tests how requireRequestEq handles validation errors
	mockTester := newMockT(t)

	// Create a request with a custom validation function that always fails
	expectedReq := Request{
		Method: "GET",
		Path:   "/test",
		Validate: func(r *http.Request) error {
			return fmt.Errorf("intentional validation error")
		},
	}

	// Create an HTTP request that would match except for the validation
	req, _ := http.NewRequest("GET", "/test", nil)

	// Call requireRequestEq directly
	requireRequestEq(mockTester, expectedReq, req)

	// Verify the test failed due to validation
	assert.True(t, mockTester.failed)
	assert.Contains(t, mockTester.errors[0], "custom validation failed")
}

// Test handling of unexpected requests
func TestServer_UnexpectedRequest(t *testing.T) {
	mockTester := newMockT(t)

	// Create a server with no expectations
	server := NewServer(mockTester, []Exchange{})
	defer server.Close()

	// Send a request that isn't expected
	resp, err := http.Get(server.Path("/unexpected"))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should have responded with 500 error and failed the test
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.True(t, mockTester.failed)
}
