package links

import (
	"net/http"
)

// HTTPClient Receiver for real http client
type HTTPClient struct{}

// WebClient Interface to web client Get
type WebClient interface {
	Get(url string) (*http.Response, error)
}

// Get fetches website body as request
func (h *HTTPClient) Get(url string) (*http.Response, error) {
	return h.Get(url)
}
