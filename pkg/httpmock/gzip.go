package httpmock

import (
	"bytes"
	"compress/gzip"
	"io"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

// decompressGzip is a BeforeSaveHook that decompresses gzipped response bodies
// to prevent binary data corruption when saving to YAML cassettes.
func decompressGzip(i *cassette.Interaction) error {
	if i.Response.Headers.Get("Content-Encoding") != "gzip" {
		return nil
	}

	gr, err := gzip.NewReader(bytes.NewReader([]byte(i.Response.Body)))
	if err != nil {
		// Leave body as-is if we can't read it
		return nil
	}
	defer func() { _ = gr.Close() }()

	decompressed, err := io.ReadAll(gr)
	if err != nil {
		// Leave body as-is if decompression fails
		return nil
	}

	// Update response with decompressed content
	i.Response.Body = string(decompressed)
	i.Response.Headers.Del("Content-Encoding")

	// Update Content-Length to match decompressed size
	i.Response.ContentLength = int64(len(decompressed))

	return nil
}
