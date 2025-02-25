// Package httpmock provides a simple, declarative API for testing HTTP clients.
//
// The package centers around the Server type, which wraps httptest.Server to provide
// a way to specify expected HTTP requests and their corresponding responses. This makes
// it easy to verify that your HTTP client makes the expected calls in the expected order.
//
// Basic usage:
//
//	server := httpmock.NewServer(t, []httpmock.Exchange{{
//		Request: httpmock.Request{
//			Method: "GET",
//			Path:   "/api/users",
//			Headers: map[string]string{
//				"Authorization": "Bearer token123",
//			},
//		},
//		Response: httpmock.Response{
//			StatusCode: http.StatusOK,
//			Body: map[string]interface{}{
//				"users": []string{"alice", "bob"},
//			},
//		},
//	}})
//	defer server.Close()
//
//	// Use server.Path("/api/users") as the base URL for your HTTP client...
//
// The Server will verify that requests match the expected method, path, headers,
// and body (if specified). For bodies, it supports both exact string matches and
// JSON comparison. You can also provide custom validation functions for more
// complex request validation:
//
//	Request: httpmock.Request{
//		Method: "POST",
//		Path:   "/api/users",
//		Validate: func(r *http.Request) error {
//			if r.Header.Get("Content-Length") == "0" {
//				return fmt.Errorf("empty request body")
//			}
//			return nil
//		},
//	}
//
// After your test completes, calling Close will shut down the server and verify that
// all expected requests were received:
//
//	server := httpmock.NewServer(t, []httpmock.Exchange{...})
//	defer server.Close() // will fail the test if not all expectations were met
//
// You can also verify expectations explicitly at any point using VerifyComplete:
//
//	if err := server.VerifyComplete(); err != nil {
//		t.Error("not all expected requests were made")
//	}
package httpmock
