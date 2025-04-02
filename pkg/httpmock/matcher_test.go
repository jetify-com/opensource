package httpmock

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

func TestCheckBasicProperties(t *testing.T) {
	tests := []struct {
		name        string
		request     *http.Request
		cassetteReq cassette.Request
		wantFail    bool
	}{
		{
			name: "all properties match",
			request: &http.Request{
				Method:     http.MethodPost,
				URL:        mustParseURL("https://api.example.com/test"),
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			cassetteReq: cassette.Request{
				Method:     http.MethodPost,
				URL:        "https://api.example.com/test",
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			wantFail: false,
		},
		{
			name: "method mismatch",
			request: &http.Request{
				Method: http.MethodPost,
				URL:    mustParseURL("https://api.example.com/test"),
			},
			cassetteReq: cassette.Request{
				Method: http.MethodGet,
				URL:    "https://api.example.com/test",
			},
			wantFail: true,
		},
		{
			name: "URL mismatch",
			request: &http.Request{
				Method: http.MethodPost,
				URL:    mustParseURL("https://api.example.com/test"),
			},
			cassetteReq: cassette.Request{
				Method: http.MethodPost,
				URL:    "https://api.example.com/different",
			},
			wantFail: true,
		},
		{
			name: "empty cassette fields match anything",
			request: &http.Request{
				Method:     http.MethodPost,
				URL:        mustParseURL("https://api.example.com/test"),
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			cassetteReq: cassette.Request{},
			wantFail:    false,
		},
		{
			name: "proto mismatch",
			request: &http.Request{
				Method:     http.MethodPost,
				URL:        mustParseURL("https://api.example.com/test"),
				Proto:      "HTTP/1.0",
				ProtoMajor: 1,
				ProtoMinor: 0,
			},
			cassetteReq: cassette.Request{
				Method:     http.MethodPost,
				URL:        "https://api.example.com/test",
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTester := &mockT{}
			got := checkBasicProperties(mockTester, tt.cassetteReq, tt.request)
			assert.Equal(t, !tt.wantFail, got)
			assert.Equal(t, tt.wantFail, mockTester.failed)
		})
	}
}

func TestCheckHeaders(t *testing.T) {
	tests := []struct {
		name     string
		actual   http.Header
		expected http.Header
		wantFail bool
	}{
		{
			name: "exact match",
			actual: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"*/*"},
			},
			expected: http.Header{
				"Content-Type": []string{"application/json"},
			},
			wantFail: false,
		},
		{
			name: "header value mismatch",
			actual: http.Header{
				"Content-Type": []string{"application/xml"},
			},
			expected: http.Header{
				"Content-Type": []string{"application/json"},
			},
			wantFail: true,
		},
		{
			name: "missing expected header",
			actual: http.Header{
				"Accept": []string{"*/*"},
			},
			expected: http.Header{
				"Content-Type": []string{"application/json"},
			},
			wantFail: true,
		},
		{
			name: "ignores extra headers in actual",
			actual: http.Header{
				"Content-Type": []string{"application/json"},
				"Extra":        []string{"value"},
			},
			expected: http.Header{
				"Content-Type": []string{"application/json"},
			},
			wantFail: false,
		},
		{
			name: "ignores case in header names",
			actual: http.Header{
				"content-type": []string{"application/json"},
			},
			expected: http.Header{
				"Content-Type": []string{"application/json"},
			},
			wantFail: false,
		},
		{
			name:     "empty expected headers match anything",
			actual:   http.Header{"Content-Type": []string{"application/json"}},
			expected: http.Header{},
			wantFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTester := &mockT{}
			got := checkHeaders(mockTester, tt.expected, tt.actual)
			assert.Equal(t, !tt.wantFail, got)
			assert.Equal(t, tt.wantFail, mockTester.failed)
		})
	}
}

func TestCheckBody(t *testing.T) {
	tests := []struct {
		name        string
		request     *http.Request
		cassetteReq cassette.Request
		wantFail    bool
	}{
		{
			name: "exact JSON match",
			request: &http.Request{
				Body: makeBody(`{"key": "value", "number": 42}`),
			},
			cassetteReq: cassette.Request{
				Body: `{"key": "value", "number": 42}`,
			},
			wantFail: false,
		},
		{
			name: "JSON field order doesn't matter",
			request: &http.Request{
				Body: makeBody(`{"number": 42, "key": "value"}`),
			},
			cassetteReq: cassette.Request{
				Body: `{"key": "value", "number": 42}`,
			},
			wantFail: false,
		},
		{
			name: "non-JSON exact match",
			request: &http.Request{
				Body: makeBody("plain text body"),
			},
			cassetteReq: cassette.Request{
				Body: "plain text body",
			},
			wantFail: false,
		},
		{
			name: "non-JSON mismatch",
			request: &http.Request{
				Body: makeBody("different text"),
			},
			cassetteReq: cassette.Request{
				Body: "plain text body",
			},
			wantFail: true,
		},
		{
			name: "empty cassette body matches anything",
			request: &http.Request{
				Body: makeBody("any content"),
			},
			cassetteReq: cassette.Request{
				Body: "",
			},
			wantFail: false,
		},
		{
			name: "nil request body only matches empty cassette body",
			request: &http.Request{
				Body: nil,
			},
			cassetteReq: cassette.Request{
				Body: "expected content",
			},
			wantFail: true,
		},
		{
			name: "error reading body",
			request: &http.Request{
				Body: errorReader{},
			},
			cassetteReq: cassette.Request{
				Body: "expected content",
			},
			wantFail: true,
		},
		{
			name: "invalid JSON in request when cassette has valid JSON",
			request: &http.Request{
				Body: makeBody(`invalid json`),
			},
			cassetteReq: cassette.Request{
				Body: `{"key": "value"}`,
			},
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTester := &mockT{}
			got := checkBody(mockTester, tt.cassetteReq.Body, tt.request)
			assert.Equal(t, !tt.wantFail, got)
			assert.Equal(t, tt.wantFail, mockTester.failed)
		})
	}
}

func TestCheckFormData(t *testing.T) {
	tests := []struct {
		name        string
		request     *http.Request
		cassetteReq cassette.Request
		wantFail    bool
	}{
		{
			name: "exact form match",
			request: func() *http.Request {
				form := url.Values{}
				form.Add("key1", "value1")
				form.Add("key2", "value2")
				req := &http.Request{
					Method: http.MethodPost,
					URL:    mustParseURL("https://example.com"),
					Header: http.Header{
						"Content-Type": []string{"application/x-www-form-urlencoded"},
					},
					Body: makeBody(form.Encode()),
				}
				if err := req.ParseForm(); err != nil {
					panic(err)
				}
				return req
			}(),
			cassetteReq: cassette.Request{
				Form: url.Values{
					"key1": []string{"value1"},
					"key2": []string{"value2"},
				},
			},
			wantFail: false,
		},
		{
			name: "form value mismatch",
			request: func() *http.Request {
				form := url.Values{}
				form.Add("key1", "wrong")
				req := &http.Request{
					Method: http.MethodPost,
					URL:    mustParseURL("https://example.com"),
					Header: http.Header{
						"Content-Type": []string{"application/x-www-form-urlencoded"},
					},
					Body: makeBody(form.Encode()),
				}
				if err := req.ParseForm(); err != nil {
					panic(err)
				}
				return req
			}(),
			cassetteReq: cassette.Request{
				Form: url.Values{
					"key1": []string{"value1"},
				},
			},
			wantFail: true,
		},
		{
			name: "missing form field",
			request: func() *http.Request {
				form := url.Values{}
				form.Add("key1", "value1")
				req := &http.Request{
					Method: http.MethodPost,
					URL:    mustParseURL("https://example.com"),
					Header: http.Header{
						"Content-Type": []string{"application/x-www-form-urlencoded"},
					},
					Body: makeBody(form.Encode()),
				}
				if err := req.ParseForm(); err != nil {
					panic(err)
				}
				return req
			}(),
			cassetteReq: cassette.Request{
				Form: url.Values{
					"key1": []string{"value1"},
					"key2": []string{"value2"},
				},
			},
			wantFail: true,
		},
		{
			name: "empty cassette form matches anything",
			request: func() *http.Request {
				form := url.Values{}
				form.Add("key1", "value1")
				req := &http.Request{
					Method: http.MethodPost,
					URL:    mustParseURL("https://example.com"),
					Header: http.Header{
						"Content-Type": []string{"application/x-www-form-urlencoded"},
					},
					Body: makeBody(form.Encode()),
				}
				if err := req.ParseForm(); err != nil {
					panic(err)
				}
				return req
			}(),
			cassetteReq: cassette.Request{
				Form: url.Values{},
			},
			wantFail: false,
		},
		{
			name: "parse form error",
			request: &http.Request{
				Method: http.MethodPost,
				URL:    mustParseURL("https://example.com"),
				Header: http.Header{
					"Content-Type": []string{"application/x-www-form-urlencoded"},
				},
				Body: makeBody("%invalid%form%data%"),
			},
			cassetteReq: cassette.Request{
				Form: url.Values{
					"key1": []string{"value1"},
				},
			},
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTester := &mockT{}
			got := checkFormData(mockTester, tt.cassetteReq.Form, tt.request)
			assert.Equal(t, !tt.wantFail, got)
			assert.Equal(t, tt.wantFail, mockTester.failed)
		})
	}
}

func TestCheckMetadata(t *testing.T) {
	tests := []struct {
		name        string
		request     *http.Request
		cassetteReq cassette.Request
		wantFail    bool
	}{
		{
			name: "all properties match",
			request: &http.Request{
				ContentLength:    100,
				TransferEncoding: []string{"chunked"},
				Host:             "example.com",
				RemoteAddr:       "192.168.1.1:1234",
				RequestURI:       "/test?key=value",
			},
			cassetteReq: cassette.Request{
				ContentLength:    100,
				TransferEncoding: []string{"chunked"},
				Host:             "example.com",
				RemoteAddr:       "192.168.1.1:1234",
				RequestURI:       "/test?key=value",
			},
			wantFail: false,
		},
		{
			name: "content length mismatch",
			request: &http.Request{
				ContentLength: 100,
			},
			cassetteReq: cassette.Request{
				ContentLength: 200,
			},
			wantFail: true,
		},
		{
			name: "transfer encoding mismatch",
			request: &http.Request{
				TransferEncoding: []string{"chunked"},
			},
			cassetteReq: cassette.Request{
				TransferEncoding: []string{"gzip"},
			},
			wantFail: true,
		},
		{
			name: "host mismatch",
			request: &http.Request{
				Host: "example.com",
			},
			cassetteReq: cassette.Request{
				Host: "different.com",
			},
			wantFail: true,
		},
		{
			name: "empty cassette properties match anything",
			request: &http.Request{
				ContentLength:    100,
				TransferEncoding: []string{"chunked"},
				Host:             "example.com",
				RemoteAddr:       "192.168.1.1:1234",
				RequestURI:       "/test?key=value",
			},
			cassetteReq: cassette.Request{},
			wantFail:    false,
		},
		{
			name: "transfer encoding order mismatch",
			request: &http.Request{
				TransferEncoding: []string{"chunked", "gzip"},
			},
			cassetteReq: cassette.Request{
				TransferEncoding: []string{"gzip", "chunked"},
			},
			wantFail: false,
		},
		{
			name: "request uri mismatch",
			request: &http.Request{
				RequestURI: "/api/v1/users",
			},
			cassetteReq: cassette.Request{
				RequestURI: "/api/v2/users",
			},
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTester := &mockT{}
			got := checkMetadata(mockTester, tt.cassetteReq, tt.request)
			assert.Equal(t, !tt.wantFail, got)
			assert.Equal(t, tt.wantFail, mockTester.failed)
		})
	}
}

func TestCheckTrailers(t *testing.T) {
	tests := []struct {
		name        string
		request     *http.Request
		cassetteReq cassette.Request
		wantFail    bool
	}{
		{
			name: "exact trailer match",
			request: &http.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value"},
				},
			},
			cassetteReq: cassette.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value"},
				},
			},
			wantFail: false,
		},
		{
			name: "trailer value mismatch",
			request: &http.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value1"},
				},
			},
			cassetteReq: cassette.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value2"},
				},
			},
			wantFail: true,
		},
		{
			name: "missing trailer",
			request: &http.Request{
				Trailer: http.Header{},
			},
			cassetteReq: cassette.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value"},
				},
			},
			wantFail: true,
		},
		{
			name: "empty cassette trailers match anything",
			request: &http.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value"},
				},
			},
			cassetteReq: cassette.Request{
				Trailer: http.Header{},
			},
			wantFail: false,
		},
		{
			name: "trailer header mismatch",
			request: &http.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value1"},
					"X-Other":   []string{"other"},
				},
			},
			cassetteReq: cassette.Request{
				Trailer: http.Header{
					"X-Trailer": []string{"value2"},
				},
			},
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTester := &mockT{}
			got := checkTrailers(mockTester, tt.cassetteReq.Trailer, tt.request.Trailer)
			assert.Equal(t, !tt.wantFail, got)
			assert.Equal(t, tt.wantFail, mockTester.failed)
		})
	}
}

func TestRequireCassetteRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     *http.Request
		cassetteReq cassette.Request
		wantFail    bool
	}{
		{
			name: "complete match",
			request: &http.Request{
				Method:     http.MethodPost,
				URL:        mustParseURL("https://api.example.com/test"),
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
					"X-Custom":     []string{"value"},
				},
				Body:             makeBody(`{"key":"value"}`),
				ContentLength:    16,
				Host:             "api.example.com",
				TransferEncoding: []string{"chunked"},
				Trailer: http.Header{
					"X-Trailer": []string{"value"},
				},
			},
			cassetteReq: cassette.Request{
				Method:     http.MethodPost,
				URL:        "https://api.example.com/test",
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body:             `{"key":"value"}`,
				ContentLength:    16,
				Host:             "api.example.com",
				TransferEncoding: []string{"chunked"},
				Trailer: http.Header{
					"X-Trailer": []string{"value"},
				},
			},
			wantFail: false,
		},
		{
			name: "partial match with minimal cassette",
			request: &http.Request{
				Method: http.MethodPost,
				URL:    mustParseURL("https://api.example.com/test"),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
					"X-Custom":     []string{"value"},
				},
				Body: makeBody(`{"key":"value"}`),
			},
			cassetteReq: cassette.Request{
				Method: http.MethodPost,
				URL:    "https://api.example.com/test",
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantFail: false,
		},
		{
			name: "mismatch in one property fails entire match",
			request: &http.Request{
				Method: http.MethodPost,
				URL:    mustParseURL("https://api.example.com/test"),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: makeBody(`{"key":"value"}`),
			},
			cassetteReq: cassette.Request{
				Method: http.MethodGet,
				URL:    "https://api.example.com/test",
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: `{"key":"value"}`,
			},
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTester := &mockT{}
			got := requireCassetteRequest(mockTester, tt.cassetteReq, tt.request)
			assert.Equal(t, !tt.wantFail, got)
			assert.Equal(t, tt.wantFail, mockTester.failed)
		})
	}
}

// Helper functions

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func makeBody(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}

// Helper type for testing read errors
type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func (errorReader) Close() error {
	return nil
}
