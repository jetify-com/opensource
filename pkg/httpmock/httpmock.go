package httpmock

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"

	"github.com/stretchr/testify/require"
)

// Request specifies what an incoming HTTP request should look like.
type Request struct {
	Method  string
	Path    string
	Body    any               // optional: if non-empty, we check the body as JSON (or raw string)
	Headers map[string]string // optional: headers that must be present
	// Validate is an optional callback for additional checks.
	Validate func(r *http.Request) error
}

// Response specifies how to respond to a matched request.
type Response struct {
	StatusCode int               // HTTP status code (defaults to 200 OK if not set)
	Body       any               // can be a JSON-marshalable object or a string
	Headers    map[string]string // optional: headers to include in response
}

// Exchange represents a pair of expected HTTP request and corresponding response.
type Exchange struct {
	Request  Request
	Response Response
}

// T is an interface that captures the testing.T methods we need
type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

// Server provides a declarative API on top of httptest.Server for testing HTTP clients.
// It allows you to specify a sequence of expected requests and their corresponding responses,
// making it easy to verify that your HTTP client makes the expected calls in the expected order.
type Server struct {
	t            T
	expectations []Exchange
	mu           sync.Mutex
	index        int
	server       *httptest.Server // make private
}

// NewServer creates and starts an httptest.Server that will match incoming requests
// to the provided expectations in order.
func NewServer(t T, expectations []Exchange) *Server {
	ds := &Server{
		t:            t,
		expectations: expectations,
		index:        0,
	}
	ds.server = httptest.NewServer(http.HandlerFunc(ds.handler))
	return ds
}

// BaseURL returns the base URL of the test server.
func (s *Server) BaseURL() string {
	return s.server.URL
}

// Path returns the complete URL for the given path.
// For example, if the server URL is "http://localhost:12345" and path is "/api/users",
// this returns "http://localhost:12345/api/users".
// url.JoinPath handles all path normalization, including handling of leading/trailing slashes.
func (s *Server) Path(p string) string {
	result, err := url.JoinPath(s.BaseURL(), p)
	require.NoError(s.t, err, "failed to join URL paths")
	return result
}

// handler checks the incoming request against the next expectation.
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.index >= len(s.expectations) {
		// No more expectations; this is an unexpected request.
		http.Error(w, "unexpected request", http.StatusInternalServerError)
		require.Fail(s.t, "Unexpected request received", "Method: %s, Path: %s", r.Method, r.URL.Path)
		return
	}

	exp := s.expectations[s.index]

	// Check if this is the expected path and method before doing detailed validation
	if exp.Request.Path != r.URL.Path || exp.Request.Method != r.Method {
		http.Error(w, "unexpected request", http.StatusInternalServerError)
		require.Fail(s.t, "Unexpected request received", "Method: %s, Path: %s", r.Method, r.URL.Path)
		return
	}

	s.index++ // move to next expectation for subsequent requests

	requireRequestEq(s.t, exp.Request, r)

	if err := writeResponse(w, exp.Response); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		require.NoError(s.t, err, "failed to write response")
	}
}

// Close shuts down the underlying httptest.Server and verifies all expected exchanges were completed.
func (s *Server) Close() {
	defer s.server.Close()
	err := s.VerifyComplete()
	require.NoError(s.t, err, "not all expectations were met")
}

// VerifyComplete checks that all expected exchanges were completed.
func (s *Server) VerifyComplete() error {
	if s.index != len(s.expectations) {
		return errors.New("not all expectations were met")
	}
	return nil
}

// requireRequestEq verifies that the actual HTTP request matches the expected request.
func requireRequestEq(tester T, expected Request, actual *http.Request) {
	// Check HTTP method.
	require.Equal(tester, expected.Method, actual.Method, "HTTP method mismatch")
	// Check path.
	require.Equal(tester, expected.Path, actual.URL.Path, "URL path mismatch")

	// Check headers if provided.
	for key, expectedValue := range expected.Headers {
		require.Equal(tester, expectedValue, actual.Header.Get(key), "Header %s mismatch", key)
	}

	// Check body if provided.
	if expected.Body != nil {
		bodyBytes, err := io.ReadAll(actual.Body)
		require.NoError(tester, err, "error reading request body")
		requireBodyEq(tester, expected.Body, bodyBytes)
	}

	// Run additional validation if provided.
	if expected.Validate != nil {
		err := expected.Validate(actual)
		require.NoError(tester, err, "custom validation failed")
	}
}

// writeResponse writes the expected response to the response writer.
func writeResponse(w http.ResponseWriter, response Response) error {
	// Set Content-Type default if not overridden
	if _, hasContentType := response.Headers["Content-Type"]; !hasContentType {
		w.Header().Set("Content-Type", "application/json")
	}

	// Set custom headers
	for k, v := range response.Headers {
		w.Header().Set(k, v)
	}

	code := response.StatusCode
	if code == 0 {
		code = http.StatusOK
	}
	w.WriteHeader(code)

	if response.Body == nil {
		return nil
	}

	// If Body is a string, write it directly. Otherwise, assume it's JSON-marshalable.
	switch rb := response.Body.(type) {
	case string:
		_, err := w.Write([]byte(rb))
		if err != nil {
			return err
		}
	default:
		if err := json.NewEncoder(w).Encode(rb); err != nil {
			return err
		}
	}
	return nil
}

// requireBodyEq verifies that the actual body matches the expected body.
// The expected body can be either a string (for exact matches) or any JSON-marshalable type.
func requireBodyEq(tester T, expected any, actualBytes []byte) {
	if expected == nil {
		return
	}

	switch expected := expected.(type) {
	case string:
		// Try JSON comparison first if it's valid JSON
		if json.Valid([]byte(expected)) {
			require.JSONEq(tester, expected, string(actualBytes))
			return
		}
		// Not valid JSON, fall back to string comparison
		require.Equal(tester, expected, string(actualBytes))
	default:
		expectedJSON, err := json.Marshal(expected)
		require.NoError(tester, err)
		require.JSONEq(tester, string(expectedJSON), string(actualBytes))
	}
}

// MergeRequests creates a new Request by combining multiple requests.
// Later requests override values from earlier requests.
// Non-zero/non-empty values from each request override values from previous requests.
func MergeRequests(requests ...Request) Request {
	if len(requests) == 0 {
		return Request{}
	}

	result := requests[0]

	for _, other := range requests[1:] {
		// Override non-empty values
		if other.Method != "" {
			result.Method = other.Method
		}
		if other.Path != "" {
			result.Path = other.Path
		}
		if other.Body != nil {
			result.Body = other.Body
		}
		if other.Validate != nil {
			result.Validate = other.Validate
		}

		// Merge headers
		if len(other.Headers) > 0 {
			if result.Headers == nil {
				result.Headers = make(map[string]string)
			}
			for k, v := range other.Headers {
				result.Headers[k] = v
			}
		}
	}

	return result
}
