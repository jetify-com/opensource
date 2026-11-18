package httpmock

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

func TestDecompressGzip(t *testing.T) {
	tests := []struct {
		name                    string
		interaction             *cassette.Interaction
		expectedBody            string
		expectedContentLength   int64
		expectedContentEncoding string
		expectNoChange          bool
	}{
		{
			name: "gzipped response is decompressed",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: createGzippedBody(t, "Hello, World!"),
					Headers: http.Header{
						"Content-Encoding": []string{"gzip"},
						"Content-Length":   []string{"32"}, // Will be updated
					},
					ContentLength: 32,
				},
			},
			expectedBody:            "Hello, World!",
			expectedContentLength:   13,
			expectedContentEncoding: "",
		},
		{
			name: "non-gzipped response unchanged",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: "Plain text response",
					Headers: http.Header{
						"Content-Type": []string{"text/plain"},
					},
					ContentLength: 19,
				},
			},
			expectedBody:            "Plain text response",
			expectedContentLength:   19,
			expectedContentEncoding: "",
			expectNoChange:          true,
		},
		{
			name: "response without content-encoding unchanged",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: "Another plain response",
					Headers: http.Header{
						"Content-Type": []string{"application/json"},
					},
					ContentLength: 22,
				},
			},
			expectedBody:            "Another plain response",
			expectedContentLength:   22,
			expectedContentEncoding: "",
			expectNoChange:          true,
		},
		{
			name: "gzipped json response",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: createGzippedBody(t, `{"message":"success","data":[1,2,3]}`),
					Headers: http.Header{
						"Content-Encoding": []string{"gzip"},
						"Content-Type":     []string{"application/json"},
						"Content-Length":   []string{"50"}, // Will be updated
					},
					ContentLength: 50,
				},
			},
			expectedBody:            `{"message":"success","data":[1,2,3]}`,
			expectedContentLength:   36, // Corrected length
			expectedContentEncoding: "",
		},
		{
			name: "invalid gzip data leaves body unchanged",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: "not gzipped data",
					Headers: http.Header{
						"Content-Encoding": []string{"gzip"},
						"Content-Length":   []string{"16"},
					},
					ContentLength: 16,
				},
			},
			expectedBody:            "not gzipped data",
			expectedContentLength:   16,
			expectedContentEncoding: "gzip", // Header should remain unchanged
			expectNoChange:          true,
		},
		{
			name: "corrupted gzip data leaves body unchanged",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x00\x00\x00\xff\xff", // Truncated gzip header
					Headers: http.Header{
						"Content-Encoding": []string{"gzip"},
						"Content-Length":   []string{"15"},
					},
					ContentLength: 15,
				},
			},
			expectedBody:            "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x00\x00\x00\xff\xff",
			expectedContentLength:   15,
			expectedContentEncoding: "gzip", // Header should remain unchanged
			expectNoChange:          true,
		},
		{
			name: "empty gzipped response",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: createGzippedBody(t, ""),
					Headers: http.Header{
						"Content-Encoding": []string{"gzip"},
						"Content-Length":   []string{"20"}, // Will be updated
					},
					ContentLength: 20,
				},
			},
			expectedBody:            "",
			expectedContentLength:   0,
			expectedContentEncoding: "",
		},
		{
			name: "gzipped response with other headers preserved",
			interaction: &cassette.Interaction{
				Response: cassette.Response{
					Body: createGzippedBody(t, "Test content"),
					Headers: http.Header{
						"Content-Encoding": []string{"gzip"},
						"Content-Type":     []string{"text/plain"},
						"Cache-Control":    []string{"no-cache"},
						"X-Custom-Header":  []string{"custom-value"},
						"Content-Length":   []string{"32"}, // Will be updated
					},
					ContentLength: 32,
				},
			},
			expectedBody:            "Test content",
			expectedContentLength:   12,
			expectedContentEncoding: "",
		},
		{
			name: "case insensitive content-encoding header",
			interaction: func() *cassette.Interaction {
				// Create headers using Set to ensure canonical form
				headers := make(http.Header)
				headers.Set("content-encoding", "gzip") // This will be canonicalized to "Content-Encoding"
				headers.Set("content-length", "29")     // Will be updated

				return &cassette.Interaction{
					Response: cassette.Response{
						Body:          createGzippedBody(t, "Case test"),
						Headers:       headers,
						ContentLength: 29,
					},
				}
			}(),
			expectedBody:            "Case test",
			expectedContentLength:   9,
			expectedContentEncoding: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original values for comparison if needed
			originalBody := tt.interaction.Response.Body
			originalHeaders := make(http.Header)
			for k, v := range tt.interaction.Response.Headers {
				originalHeaders[k] = append([]string(nil), v...)
			}
			originalContentLength := tt.interaction.Response.ContentLength

			// Call the function
			err := decompressGzip(tt.interaction)

			// Function should never return an error
			assert.NoError(t, err)

			// Check the body
			assert.Equal(t, tt.expectedBody, tt.interaction.Response.Body)

			// Check content length
			assert.Equal(t, tt.expectedContentLength, tt.interaction.Response.ContentLength)

			// Check Content-Encoding header
			actualContentEncoding := tt.interaction.Response.Headers.Get("Content-Encoding")
			assert.Equal(t, tt.expectedContentEncoding, actualContentEncoding)

			if tt.expectNoChange {
				// Verify that nothing changed for non-gzip responses or invalid gzip
				assert.Equal(t, originalBody, tt.interaction.Response.Body)
				assert.Equal(t, originalContentLength, tt.interaction.Response.ContentLength)

				// For invalid gzip, the Content-Encoding header should remain
				if tt.expectedContentEncoding != "" {
					assert.Equal(t, originalHeaders.Get("Content-Encoding"), tt.interaction.Response.Headers.Get("Content-Encoding"))
				}
			} else {
				// Verify that other headers are preserved (except Content-Encoding and Content-Length)
				for key, values := range originalHeaders {
					if key != "Content-Encoding" && key != "Content-Length" {
						assert.Equal(t, values, tt.interaction.Response.Headers[key], "Header %s should be preserved", key)
					}
				}
			}
		})
	}
}

// Helper function to create gzipped content for testing
func createGzippedBody(t *testing.T, content string) string {
	t.Helper()

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)

	_, err := gw.Write([]byte(content))
	require.NoError(t, err)

	err = gw.Close()
	require.NoError(t, err)

	return buf.String()
}
