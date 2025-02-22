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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 1. Examples
func ExampleServer_basic() {
	ts := NewServer(&mockT{}, []Exchange{{
		Request: Request{
			Method: "GET",
			Path:   "/hello",
		},
		Response: Response{
			Body: "world",
		},
	}})
	defer ts.Close()

	resp, _ := http.Get(ts.Path("/hello"))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: world
}

func ExampleServer_jsonRequest() {
	ts := NewServer(&mockT{}, []Exchange{{
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
	defer ts.Close()

	resp, _ := http.Post(ts.Path("/api/users"),
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
	ts := NewServer(&mockT{}, []Exchange{
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
	defer ts.Close()

	// Login
	resp, _ := http.Post(ts.Path("/login"),
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
	req, _ := http.NewRequest("GET", ts.Path("/profile"), nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp, _ = http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: {"name":"Alice"}
}

func ExampleServer_validation() {
	ts := NewServer(&mockT{}, []Exchange{{
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
	defer ts.Close()

	// Send non-empty request
	resp, _ := http.Post(ts.Path("/upload"),
		"text/plain",
		strings.NewReader("some data"))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: uploaded
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := newMockT(t)
			ts := NewServer(mt, []Exchange{tt.expect})
			defer ts.Close()

			reqToSend := tt.expect.Request
			if tt.send != nil {
				reqToSend = *tt.send
			}

			req, err := buildRequest(ts.BaseURL(), reqToSend)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if tt.wantFail {
				assert.True(t, mt.failed)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			} else {
				assert.False(t, mt.failed)
				assertResponseEq(t, tt.expect.Response, resp)
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
			wantErr: "not all expectations were met",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := newMockT(t)
			ts := NewServer(mt, tt.expect)
			defer ts.Close()

			// Send all requests in order
			for _, req := range tt.send {
				r, err := buildRequest(ts.BaseURL(), req)
				require.NoError(t, err)
				resp, err := http.DefaultClient.Do(r)
				require.NoError(t, err)
				resp.Body.Close()
			}

			err := ts.VerifyComplete()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := newMockT(t)
			mt.failed = false
			mt.errors = nil

			requireBodyEq(mt, tt.expected, []byte(tt.actual))
			assert.Equal(t, tt.wantFail, mt.failed)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a server with a custom base URL for testing
			s := &Server{server: httptest.NewServer(nil)}
			s.server.URL = tt.baseURL // override the random port with our test URL
			defer s.server.Close()

			got := s.Path(tt.path)
			assert.Equal(t, tt.expected, got)
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
func buildRequest(baseURL string, r Request) (*http.Request, error) {
	u, err := url.JoinPath(baseURL, r.Path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, u, buildBody(r.Body))
	if err != nil {
		return nil, err
	}

	// Set headers from the request
	for k, v := range r.Headers {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeRequests(tt.base, tt.override)
			assert.Equal(t, tt.want, got)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewServer(t, []Exchange{{
				Request: Request{
					Method: "GET",
					Path:   "/test",
				},
				Response: tt.response,
			}})
			defer ts.Close()

			resp, err := http.Get(ts.Path("/test"))
			require.NoError(t, err)
			defer resp.Body.Close()

			for k, want := range tt.wantHeaders {
				assert.Equal(t, want, resp.Header.Get(k))
			}
		})
	}
}
