package httpmock_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"go.jetify.com/pkg/httpmock"
)

// t is a simple implementation that satisfies the testing.TB interface
// requirements for the examples, but is separate from the mockT in the test file.
type t struct{}

func (*t) Errorf(format string, args ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}
func (*t) FailNow() {}
func (*t) Helper()  {}

// ExampleServer_basic demonstrates basic usage of the httpmock Server.
func ExampleServer_basic() {
	testServer := httpmock.NewServer(&t{}, []httpmock.Exchange{{
		Request: httpmock.Request{
			Method: "GET",
			Path:   "/hello",
		},
		Response: httpmock.Response{
			Body: "world",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		},
	}})
	defer testServer.Close()

	resp, _ := http.Get(testServer.Path("/hello"))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: world
}

// ExampleServer_jsonRequest demonstrates handling JSON requests and responses.
func ExampleServer_jsonRequest() {
	testServer := httpmock.NewServer(&t{}, []httpmock.Exchange{{
		Request: httpmock.Request{
			Method: "POST",
			Path:   "/api/users",
			Body:   `{"name":"Alice"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		Response: httpmock.Response{
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

// ExampleServer_sequence demonstrates a sequence of requests and responses.
func ExampleServer_sequence() {
	testServer := httpmock.NewServer(&t{}, []httpmock.Exchange{
		{
			Request: httpmock.Request{
				Method: "POST",
				Path:   "/login",
				Body:   `{"username":"alice","password":"secret"}`,
			},
			Response: httpmock.Response{
				Body: map[string]string{"token": "abc123"},
			},
		},
		{
			Request: httpmock.Request{
				Method: "GET",
				Path:   "/profile",
				Headers: map[string]string{
					"Authorization": "Bearer abc123",
				},
			},
			Response: httpmock.Response{
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
	_ = resp.Body.Close()

	// Get profile using token
	req, _ := http.NewRequest(http.MethodGet, testServer.Path("/profile"), nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp, _ = http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output: {"name":"Alice"}
}

// ExampleServer_validation demonstrates custom request validation.
func ExampleServer_validation() {
	testServer := httpmock.NewServer(&t{}, []httpmock.Exchange{{
		Request: httpmock.Request{
			Method: "POST",
			Path:   "/upload",
			Validate: func(r *http.Request) error {
				if r.Header.Get("Content-Length") == "0" {
					return fmt.Errorf("empty request body")
				}
				return nil
			},
		},
		Response: httpmock.Response{
			StatusCode: http.StatusOK,
			Body:       "uploaded",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
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

// ExampleServer_delay demonstrates response delay functionality.
func ExampleServer_delay() {
	// Create a server with a delayed response
	testServer := httpmock.NewServer(&t{}, []httpmock.Exchange{{
		Request: httpmock.Request{
			Method: "GET",
			Path:   "/api/slow",
		},
		Response: httpmock.Response{
			StatusCode: http.StatusOK,
			Body:       `{"status":"success", "data":"worth the wait"}`,
			Delay:      1 * time.Second, // Response will be delayed by 1 second
		},
	}})
	defer testServer.Close()

	// Print the configured delay
	fmt.Println("Configured delay:", 1*time.Second)

	// Make the request
	resp, err := http.Get(testServer.Path("/api/slow"))
	if err != nil {
		fmt.Println("request error:", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	// Output:
	// Configured delay: 1s
}

// ExampleReplayServer_basic demonstrates basic usage of the ReplayServer for recording
// and replaying HTTP interactions with httpbin.org.
func ExampleReplayServer_basic() {
	mockT := &t{} // In a real test, this would be the actual testing.T instance

	// Create a new replay server that will record interactions with httpbin.org
	replayServer, err := httpmock.NewReplayServer(mockT, httpmock.ReplayConfig{
		Host:     "https://httpbin.org",
		Cassette: "testdata/successful_get",
	})
	require.NoError(mockT, err)
	defer func() { _ = replayServer.Close() }()

	// Make a request that will be recorded or replayed
	req, err := http.NewRequest(http.MethodGet, replayServer.URL()+"/get", nil)
	require.NoError(mockT, err)
	req.Host = "httpbin.org"
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(mockT, err)
	defer func() { _ = resp.Body.Close() }()

	// Read and print the response status and body
	body, err := io.ReadAll(resp.Body)
	require.NoError(mockT, err)

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(mockT, err)

	// Print only the status code
	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}

// ExampleReplayServer_headers demonstrates how ReplayServer handles headers
// while recording and replaying HTTP interactions.
func ExampleReplayServer_headers() {
	mockT := &t{} // In a real test, this would be the actual testing.T instance

	replayServer, err := httpmock.NewReplayServer(mockT, httpmock.ReplayConfig{
		Host:     "https://httpbin.org",
		Cassette: "testdata/get_with_header",
	})
	require.NoError(mockT, err)
	defer func() { _ = replayServer.Close() }()

	// Create a request with custom headers
	req, err := http.NewRequest(http.MethodGet, replayServer.URL()+"/headers", nil)
	require.NoError(mockT, err)
	req.Host = "httpbin.org"
	req.Header.Set("X-Test-Header", "test-value")
	req.Header.Set("Accept-Encoding", "gzip")

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(mockT, err)
	defer func() { _ = resp.Body.Close() }()

	// Print response status
	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}

// ExampleReplayServer_jsonRequest demonstrates handling JSON requests with ReplayServer.
func ExampleReplayServer_jsonRequest() {
	mockT := &t{} // In a real test, this would be the actual testing.T instance

	replayServer, err := httpmock.NewReplayServer(mockT, httpmock.ReplayConfig{
		Host:     "https://httpbin.org",
		Cassette: "testdata/post_with_body",
	})
	require.NoError(mockT, err)
	defer func() { _ = replayServer.Close() }()

	// Create a JSON request
	jsonData := `{"test": "data"}`
	req, err := http.NewRequest(http.MethodPost, replayServer.URL()+"/post", strings.NewReader(jsonData))
	require.NoError(mockT, err)
	req.Host = "httpbin.org"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(jsonData)))

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(mockT, err)
	defer func() { _ = resp.Body.Close() }()

	// Print response status
	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 200
}
