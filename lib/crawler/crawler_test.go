package crawler_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smashed-avo/go-crawler/lib/crawler"
	"github.com/smashed-avo/go-crawler/lib/data"
)

const (
	_ mockStateWorker = iota
	emptyResponse
	successResponse
	maxDepthReachedResponse
)

var (
	parent = &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: make([]*data.Response, 0)}
	child1 = &data.Response{Depth: 1, Title: "", URL: "https://www.successweb.com/children1", Nodes: make([]*data.Response, 0)}
	child2 = &data.Response{Depth: 2, Title: "", URL: "https://www.successweb.com/children2", Nodes: make([]*data.Response, 0)}
	child3 = &data.Response{Depth: 3, Title: "", URL: "https://www.successweb.com/children3", Nodes: make([]*data.Response, 0)}
	child4 = &data.Response{Depth: 4, Title: "", URL: "https://www.successweb.com/children4", Nodes: make([]*data.Response, 0)}
)

type mockStateWorker int

type MockWorker struct {
	State   mockStateWorker
	ChQueue chan []*data.Response
}

func (w *MockWorker) Do(node *data.Response, depth int, chQueue chan []*data.Response, visited *data.Visited) {
	switch w.State {
	case emptyResponse:
		nodes := make([]*data.Response, 0)
		node.Nodes = nodes
		chQueue <- nodes
		return
	case successResponse, maxDepthReachedResponse:
		nodes := make([]*data.Response, 0)
		switch depth {
		case 1:
			nodes = append(nodes, child1)
			break
		case 2:
			nodes = append(nodes, child2)
			break
		case 3:
			nodes = append(nodes, child3)
			break
		case 4:
			nodes = append(nodes, child4)
			break
		}
		node.Nodes = nodes
		chQueue <- nodes
		return
	default:
		panic(fmt.Sprintf("Invalid mockStateWorker: %v", w.State))
	}
}

func (w *MockWorker) GetPageTitle(u string) string {
	switch w.State {
	case emptyResponse, successResponse, maxDepthReachedResponse:
		return "Success Web"
	default:
		panic(fmt.Sprintf("Invalid mockStateWorker: %v", w.State))
	}
}

func TestCrawl(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		name             string
		state            mockStateWorker
		maxDepth         int
		url              string
		title            string
		expectedResponse *data.Response
	}{
		{
			name:     "Success",
			state:    successResponse,
			maxDepth: 0,
			url:      "https://www.successweb.com",
			title:    "Success Web",
			expectedResponse: &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{
				&data.Response{Depth: 1, Title: "", URL: "https://www.successweb.com/children1", Nodes: []*data.Response{}}}},
		},
		{
			name:             "Empty",
			state:            emptyResponse,
			maxDepth:         0,
			url:              "https://www.successweb.com",
			title:            "Success Web",
			expectedResponse: &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{}},
		},
		{
			name:     "Reach Max Depth 2",
			state:    maxDepthReachedResponse,
			maxDepth: 2,
			url:      "https://www.successweb.com",
			title:    "Success Web",
			expectedResponse: &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{
				&data.Response{Depth: 1, Title: "", URL: "https://www.successweb.com/children1", Nodes: []*data.Response{}}}},
		},
		{
			name:     "Reach Max Depth 3",
			state:    maxDepthReachedResponse,
			maxDepth: 3,
			url:      "https://www.successweb.com",
			title:    "Success Web",
			expectedResponse: &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{
				&data.Response{Depth: 1, Title: "", URL: "https://www.successweb.com/children1", Nodes: []*data.Response{
					&data.Response{Depth: 2, Title: "", URL: "https://www.successweb.com/children2", Nodes: []*data.Response{}}}}}},
		},
		{
			name:     "Reach Max Depth 4",
			state:    maxDepthReachedResponse,
			maxDepth: 4,
			url:      "https://www.successweb.com",
			title:    "Success Web",
			expectedResponse: &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{
				&data.Response{Depth: 1, Title: "", URL: "https://www.successweb.com/children1", Nodes: []*data.Response{
					&data.Response{Depth: 2, Title: "", URL: "https://www.successweb.com/children2", Nodes: []*data.Response{
						&data.Response{Depth: 3, Title: "", URL: "https://www.successweb.com/children3", Nodes: []*data.Response{}}}}}}}},
		},
		{
			name:     "Reach Max Depth 5",
			state:    maxDepthReachedResponse,
			maxDepth: 5,
			url:      "https://www.successweb.com",
			title:    "Success Web",
			expectedResponse: &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{
				&data.Response{Depth: 1, Title: "", URL: "https://www.successweb.com/children1", Nodes: []*data.Response{
					&data.Response{Depth: 2, Title: "", URL: "https://www.successweb.com/children2", Nodes: []*data.Response{
						&data.Response{Depth: 3, Title: "", URL: "https://www.successweb.com/children3", Nodes: []*data.Response{
							&data.Response{Depth: 4, Title: "", URL: "https://www.successweb.com/children4", Nodes: []*data.Response{}}}}}}}}}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.ParseRequestURI(tc.url)
			assert.NoError(err)

			m := MockWorker{State: tc.state}
			c := crawler.NewCrawler(&m)

			c.ChQueue = make(chan []*data.Response)
			defer close(c.ChQueue)

			r := c.Crawl(u, tc.maxDepth)

			assert.Equal(tc.expectedResponse, r, tc.name)
		})
	}
}
