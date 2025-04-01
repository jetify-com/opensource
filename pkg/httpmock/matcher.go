package httpmock

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

func requireCassetteRequest(tester T, expected cassette.Request, actual *http.Request) bool {
	tester.Helper()

	// Check each component of the request
	return checkBasicProperties(tester, expected, actual) &&
		checkHeaders(tester, expected.Headers, actual.Header) &&
		checkBody(tester, expected.Body, actual) &&
		checkFormData(tester, expected.Form, actual) &&
		checkMetadata(tester, expected, actual) &&
		checkTrailers(tester, expected.Trailer, actual.Trailer)
}

// checkBasicProperties validates the basic HTTP request properties like method, URL, and protocol
func checkBasicProperties(tester T, expected cassette.Request, actual *http.Request) bool {
	return assert.Condition(tester, func() bool {
		if expected.Method != "" {
			if !assert.Equal(tester, expected.Method, actual.Method, "method mismatch") {
				return false
			}
		}
		if expected.URL != "" {
			if !assert.Equal(tester, expected.URL, actual.URL.String(), "URL mismatch") {
				return false
			}
		}
		if expected.Proto != "" {
			if !assert.Equal(tester, expected.Proto, actual.Proto, "protocol mismatch") {
				return false
			}
		}
		if expected.ProtoMajor != 0 {
			if !assert.Equal(tester, expected.ProtoMajor, actual.ProtoMajor, "protocol major version mismatch") {
				return false
			}
		}
		if expected.ProtoMinor != 0 {
			if !assert.Equal(tester, expected.ProtoMinor, actual.ProtoMinor, "protocol minor version mismatch") {
				return false
			}
		}
		return true
	})
}

// checkHeaders validates the HTTP headers match the expected values
func checkHeaders(tester T, expected http.Header, actual http.Header) bool {
	filteredActual := make(http.Header)
	for k, v := range actual {
		filteredActual[http.CanonicalHeaderKey(k)] = v
	}
	for _, header := range ignoredHeaders {
		filteredActual.Del(header)
	}

	return assert.Condition(tester, func() bool {
		for k, v := range expected {
			canonicalKey := http.CanonicalHeaderKey(k)
			actualValues := filteredActual[canonicalKey]
			if !assert.NotEmpty(tester, actualValues, "missing header %q", k) {
				return false
			}
			if !assert.ElementsMatch(tester, v, actualValues, "header %q values mismatch", k) {
				return false
			}
		}
		return true
	})
}

// checkBody validates the request body matches the expected value
func checkBody(tester T, expected string, actual *http.Request) bool {
	if expected == "" {
		return true
	}

	if !assert.NotNil(tester, actual.Body, "missing request body") {
		return false
	}

	bodyBytes, err := io.ReadAll(actual.Body)
	if !assert.NoError(tester, err, "error reading request body") {
		return false
	}
	actual.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Try JSON comparison first
	if isJSON(expected) && isJSON(string(bodyBytes)) {
		return assert.JSONEq(tester, expected, string(bodyBytes), "JSON body mismatch")
	}

	return assert.Equal(tester, expected, string(bodyBytes), "body mismatch")
}

// isJSON checks if a string is valid JSON
func isJSON(str string) bool {
	if str == "" {
		return false
	}
	str = strings.TrimSpace(str)
	if str == "" || (str[0] != '{' && str[0] != '[') {
		return false
	}
	var js interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

// checkFormData validates the form data matches the expected values
func checkFormData(tester T, expected map[string][]string, actual *http.Request) bool {
	if len(expected) == 0 {
		return true
	}

	if !assert.NoError(tester, actual.ParseForm(), "error parsing form data") {
		return false
	}

	return assert.Condition(tester, func() bool {
		for k, expectedValues := range expected {
			actualValues := actual.Form[k]
			if !assert.NotEmpty(tester, actualValues, "missing form field %q", k) {
				return false
			}
			if !assert.ElementsMatch(tester, expectedValues, actualValues, "form field %q values mismatch", k) {
				return false
			}
		}
		return true
	})
}

// checkMetadata validates request metadata like content length and transfer encoding
func checkMetadata(tester T, expected cassette.Request, actual *http.Request) bool {
	return assert.Condition(tester, func() bool {
		if expected.ContentLength != 0 {
			if !assert.Equal(tester, expected.ContentLength, actual.ContentLength, "content length mismatch") {
				return false
			}
		}
		if len(expected.TransferEncoding) > 0 {
			if !assert.ElementsMatch(tester, expected.TransferEncoding, actual.TransferEncoding, "transfer encoding mismatch") {
				return false
			}
		}
		if expected.Host != "" {
			if !assert.Equal(tester, expected.Host, actual.Host, "host mismatch") {
				return false
			}
		}
		if expected.RequestURI != "" {
			if !assert.Equal(tester, expected.RequestURI, actual.RequestURI, "request URI mismatch") {
				return false
			}
		}
		return true
	})
}

// checkTrailers validates the trailer headers match the expected values
func checkTrailers(tester T, expected http.Header, actual http.Header) bool {
	if len(expected) == 0 {
		return true
	}

	filteredTrailer := make(http.Header)
	for k, v := range actual {
		filteredTrailer[http.CanonicalHeaderKey(k)] = v
	}
	for _, header := range ignoredHeaders {
		filteredTrailer.Del(header)
	}

	return assert.Condition(tester, func() bool {
		for k, v := range expected {
			canonicalKey := http.CanonicalHeaderKey(k)
			actualValues := filteredTrailer[canonicalKey]
			if !assert.NotEmpty(tester, actualValues, "missing trailer header %q", k) {
				return false
			}
			if !assert.ElementsMatch(tester, v, actualValues, "trailer header %q values mismatch", k) {
				return false
			}
		}
		return true
	})
}
