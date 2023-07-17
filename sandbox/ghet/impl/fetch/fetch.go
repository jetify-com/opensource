package fetch

import "net/http"

func (c *Client) Fetch(url string) (*http.Response, error) {
	return c.httpClient.Get(url)
}

func Fetch(url string) (*http.Response, error) {
	return defaultClient.Fetch(url)
}
