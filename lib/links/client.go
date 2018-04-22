package links

import (
	"io"
	"net/http"
	"time"
)

// HTTPClient Receiver for real http client
type HTTPClient struct {
}

// WebClient Interface to web client Get
type WebClient interface {
	Get(url string) (io.Reader, error)
}

// Get fetches website body as request
func (c *HTTPClient) Get(url string) (io.Reader, error) {
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Body, nil
}
