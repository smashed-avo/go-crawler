package webpage_test

import (
	"errors"
	"testing"
	"net/http"
	"io"
	"bytes"

	"github.com/stretchr/testify/assert"

	"github.com/smashed-avo/go-crawler/lib/webpage"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type successClient struct{
	body string
}

func (c successClient) Get(_ string) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: nopCloser{bytes.NewBufferString(c.body)}}, nil
}

type errorClient struct{}

func (c errorClient) Get(_ string) (*http.Response, error) {
	return nil, errors.New("Error communicating with server")
}

func TestLinkFetcher(t *testing.T) {
	tt := []struct {
		name            string
		body            string
		stubClient      *http.Client
		expectedResult  string
	}{
		{
			name:            "Success POST",
			body:            "{}",
			// TODO Mock http client
			stubClient:      &successClient{},
			expectedResult:  "",
		},
		{
			name:            "Client FAILED",
			body:            "{}",
			stubClient:      &errorClient{},
			expectedResult:  "",
		},
	}

	for _, tc := range tt {
		mockLinks := make(chan string)
		mockFinished := make(chan bool)
		mockErrors := make(chan error)

		fetcher := webpage.LinkFetcher{tc.stubClient}

		go fetcher.Links("www.test.com", mockLinks, mockFinished, mockErrors)

		// Assert results and errors properly
		result := <-mockLinks
		finished := <-mockFinished
		errors := <-mockErrors

		assert.Equal(t, tc.expectedResult, result, "Result mismatch")

		close(mockLinks)
		close(mockFinished)
		close(mockErrors)

	}

}

