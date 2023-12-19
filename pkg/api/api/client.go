package api

import (
	"net/http"

	"go.jetpack.io/pkg/auth/session"
)

// Client manages state for interacting with the JetCloud API, as well as
// communicating with the JetCloud API.
type Client struct {
	Host       string
	httpClient *http.Client
}

func NewClient(host string, token *session.Token) *Client {
	return &Client{
		Host:       host,
		httpClient: &http.Client{Transport: &transport{token}},
	}
}

type transport struct{ token *session.Token }

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.token.AccessToken)
	return http.DefaultTransport.RoundTrip(req)
}
