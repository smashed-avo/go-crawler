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
		return &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: make([]*data.Response, 0)}
	default:
		panic(fmt.Sprintf("Invalid mockStateCrawler: %v", c.State))
	}
}

// GET /crawl
func TestHandleCrawl(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		name               string
		state              mockStateCrawler
		url                string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Get success",
			state:              successResponse,
			url:                "/crawl?url=https://www.successweb.com",
			expectedStatusCode: 200,
			expectedBody:       `{"depth":0,"title":"Success Web","url":"https://www.successweb.com","nodes":[]}`,
		},
		{
			name:               "Success: Empty depth defaulted",
			state:              successResponse,
			url:                "/crawl?url=https://successweb.com&depth=",
			expectedStatusCode: 200,
			expectedBody:       `{"depth":0,"title":"Success Web","url":"https://www.successweb.com","nodes":[]}`,
		},
		{
			name:               "Success: passing depth",
			state:              successResponse,
			url:                "/crawl?url=https://successweb.com?depth=5",
			expectedStatusCode: 200,
			expectedBody:       `{"depth":0,"title":"Success Web","url":"https://www.successweb.com","nodes":[]}`,
		},
		{
			name:               "Bad Request: empty URL",
			state:              emptyResponse,
			url:                "/crawl",
			expectedStatusCode: 400,
			expectedBody:       ``,
		},
		{
			name:               "Bad Request: non parseable URL",
			state:              emptyResponse,
			url:                "/crawl?url=http//notanurl.com",
			expectedStatusCode: 400,
			expectedBody:       ``,
		},
		{
			name:               "Bad Request: depth not int",
			state:              emptyResponse,
			url:                "/crawl?url=https://www.successweb.com&depth=aaaa",
			expectedStatusCode: 400,
			expectedBody:       ``,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h := handler.NewHandler(&MockCrawler{State: tc.state})

			req, err := http.NewRequest("GET", tc.url, nil)
			assert.NoError(err)

			w := httptest.NewRecorder()
			h.HandleCrawl(w, req)

			assert.Equal(tc.expectedStatusCode, w.Code, tc.name)

			body, err := ioutil.ReadAll(w.Body)
			require.NoError(t, err, "Error reading response")
			if string(body) != "" {
				assert.JSONEq(tc.expectedBody, string(body), tc.name)
			} else {
				assert.Equal(tc.expectedBody, string(body), tc.name)
			}
		})
	}
}
