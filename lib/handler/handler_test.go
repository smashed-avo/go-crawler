package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smashed-avo/go-crawler/lib/data"
	"github.com/smashed-avo/go-crawler/lib/handler"
)

const (
	_ mockStateCrawler = iota
	emptyResponse
	successResponse
)

type mockStateCrawler int

type MockCrawler struct {
	State mockStateCrawler
}

func (c *MockCrawler) Crawl(seedURL *url.URL, maxDepth int) *data.Response {
	switch c.State {
	case emptyResponse:
		return &data.Response{}
	case successResponse:
		return &data.Response{}
	default:
		panic(fmt.Sprintf("Invalid mockStateCrawler: %v", c.State))
	}
}

// GET /crawl
func TestHandleCrawl(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		Name               string
		state              mockStateCrawler
		url                string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			Name:               "Get success",
			state:              successResponse,
			url:                "/crawl?url=hppts://successweb.com",
			expectedStatusCode: 200,
			expectedBody:       `{"depth":0,"title":"","url":"","nodes":null}`,
		},
		{
			Name:               "Bad Request: empty URL",
			state:              emptyResponse,
			url:                "/crawl",
			expectedStatusCode: 400,
			expectedBody:       ``,
		},
		{
			Name:               "Success: Empty depth defaulted",
			state:              emptyResponse,
			url:                "/crawl?url=hppts://successweb.com?depth=",
			expectedStatusCode: 200,
			expectedBody:       `{"depth":0,"title":"","url":"","nodes":null}`,
		},
		{
			Name:               "Success: depth not int defaulted",
			state:              emptyResponse,
			url:                "/crawl?url=hppts://successweb.com?depth=12a",
			expectedStatusCode: 200,
			expectedBody:       `{"depth":0,"title":"","url":"","nodes":null}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			h := &handler.Handler{Crawler: &MockCrawler{State: tc.state}}

			req, err := http.NewRequest("GET", tc.url, nil)
			assert.NoError(err)

			w := httptest.NewRecorder()
			h.HandleCrawl(w, req)

			assert.Equal(tc.expectedStatusCode, w.Code, tc.Name)

			body, err := ioutil.ReadAll(w.Body)
			require.NoError(t, err, "Error reading response")
			if string(body) != "" {
				assert.JSONEq(tc.expectedBody, string(body), tc.Name)
			} else {
				assert.Equal(tc.expectedBody, string(body), tc.Name)
			}
		})
	}
}