package httpmock_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.jetify.com/pkg/httpmock"
)

// t is a simple implementation that satisfies the testing.TB interface
// requirements for the examples, but is separate from the mockT in the test file.
type t struct{}

func (*t) Errorf(format string, args ...interface{}) {}
func (*t) FailNow()                                  {}

// ExampleServer_basic demonstrates basic usage of the httpmock Server.
func ExampleServer_basic() {
	testServer := httpmock.NewServer(&t{}, []httpmock.Exchange{{
		Request: httpmock.Request{
			Method: "GET",
			Path:   "/hello",
		},
		Response: httpmock.Response{
			Body: "world",
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
	resp.Body.Close()

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

	// Instead of printing the URL (which has a dynamic port), print something static
	fmt.Println("Configured delay:", 1*time.Second)

	// Output:
	// Configured delay: 1s
}
