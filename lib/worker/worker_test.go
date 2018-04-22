package worker_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smashed-avo/go-crawler/lib/data"
	"github.com/smashed-avo/go-crawler/lib/worker"
)

var (
	parent            = &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: make([]*data.Response, 0)}
	child1            = &data.Response{Depth: 1, Title: "", URL: "https://www.successweb.com/children1", Nodes: make([]*data.Response, 0)}
	child2            = &data.Response{Depth: 2, Title: "", URL: "https://www.successweb.com/children2", Nodes: make([]*data.Response, 0)}
	child3            = &data.Response{Depth: 3, Title: "", URL: "https://www.successweb.com/children3", Nodes: make([]*data.Response, 0)}
	child4            = &data.Response{Depth: 4, Title: "", URL: "https://www.successweb.com/children4", Nodes: make([]*data.Response, 0)}
	link1             = "www.fakeweb.com/test1"
	link2             = "www.fakeweb.com/test2"
	link3             = "www.fakeweb.com/test3"
	linkNonParseable  = "http://a b.com/"
	link1node         = data.Response{Depth: 1, Title: "", URL: "www.fakeweb.com/test1", Nodes: []*data.Response{}}
	link2node         = data.Response{Depth: 1, Title: "", URL: "www.fakeweb.com/test2", Nodes: []*data.Response{}}
	link3node         = data.Response{Depth: 1, Title: "", URL: "www.fakeweb.com/test3", Nodes: []*data.Response{}}
	linkWithTitle     = "https://en.wikipedia.org/wiki/Go_(programming_language)"
	linkWithTitleNode = data.Response{Depth: 1, Title: "Go (programming language) - Wikipedia", URL: "https://en.wikipedia.org/wiki/Go_(programming_language)", Nodes: []*data.Response{}}
)

const (
	_ mockStateCollector = iota
	successThreeLinksFinished
	successLinkWithTitleFinished
	successRepeatedLinkFinished
	successNonParseableLinkNotIncluded
	errored
)

type mockStateCollector int

type MockCollector struct {
	State      mockStateCollector
	ChLinks    chan string
	ChErrors   chan error
	ChFinished chan bool
}

func (w *MockCollector) Collect(url string, chLinks chan string, chFinished chan bool, chErrors chan error) {
	switch w.State {
	case successThreeLinksFinished:
		chLinks <- link1
		chLinks <- link2
		chLinks <- link3
		chFinished <- true
		return
	case successLinkWithTitleFinished:
		chLinks <- linkWithTitle
		chFinished <- true
		return
	case successRepeatedLinkFinished:
		chLinks <- link1
		chLinks <- link2
		chLinks <- link3
		chLinks <- link1
		chFinished <- true
		return
	case successNonParseableLinkNotIncluded:
		chLinks <- link1
		chLinks <- link2
		chLinks <- linkNonParseable
		chLinks <- link3
		chFinished <- true
		return
	case errored:
		chErrors <- errors.New("Test error")
		return
	default:
		panic(fmt.Sprintf("Invalid mockStateCollector: %v", w.State))
	}
}

func TestDo(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		name                string
		state               mockStateCollector
		depth               int
		node                *data.Response
		expectedQueueValues []*data.Response
		expectedVisited     *data.Visited
		expectedNode        *data.Response
	}{
		{
			name:                "Success - Collect three links",
			state:               successThreeLinksFinished,
			depth:               1,
			node:                &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{}},
			expectedQueueValues: []*data.Response{&link1node, &link2node, &link3node},
			expectedVisited:     addVisited(&data.Visited{M: make(map[string]bool)}, link1, link2, link3),
			expectedNode:        &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{&link1node, &link2node, &link3node}},
		},
		{
			name:                "Success - Collect link with title",
			state:               successLinkWithTitleFinished,
			depth:               1,
			node:                &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{}},
			expectedQueueValues: []*data.Response{&linkWithTitleNode},
			expectedVisited:     addVisited(&data.Visited{M: make(map[string]bool)}, linkWithTitle),
			expectedNode:        &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{&linkWithTitleNode}},
		},
		{
			name:                "Success - Repeated link",
			state:               successRepeatedLinkFinished,
			depth:               1,
			node:                &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{}},
			expectedQueueValues: []*data.Response{&link1node, &link2node, &link3node},
			expectedVisited:     addVisited(&data.Visited{M: make(map[string]bool)}, link1, link2, link3),
			expectedNode:        &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{&link1node, &link2node, &link3node}},
		},
		{
			name:                "Success - Non parseable link excluded",
			state:               successNonParseableLinkNotIncluded,
			depth:               1,
			node:                &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{}},
			expectedQueueValues: []*data.Response{&link1node, &link2node, &link3node},
			expectedVisited:     addVisited(&data.Visited{M: make(map[string]bool)}, link1, link2, link3),
			expectedNode:        &data.Response{Depth: 0, Title: "Success Web", URL: "https://www.successweb.com", Nodes: []*data.Response{&link1node, &link2node, &link3node}},
		},
		{
			name:                "Error - couldn't connect to site",
			state:               errored,
			depth:               1,
			node:                &data.Response{Depth: 0, Title: "Error Web", URL: "https://www.errorweb.com", Nodes: []*data.Response{}},
			expectedQueueValues: []*data.Response{},
			expectedVisited:     &data.Visited{M: make(map[string]bool)},
			expectedNode:        &data.Response{Depth: 0, Title: "Error Web", URL: "https://www.errorweb.com", Nodes: []*data.Response{}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := MockCollector{State: tc.state}
			w := worker.NewWorker(&m)

			w.ChLinks = make(chan string)
			w.ChErrors = make(chan error)
			w.ChFinished = make(chan bool)
			defer close(w.ChLinks)
			defer close(w.ChErrors)
			defer close(w.ChFinished)

			q := make(chan []*data.Response)
			v := data.Visited{M: make(map[string]bool)}

			go w.Do(tc.node, tc.depth, q, &v)

			values := <-q

			assert.Equal(tc.expectedVisited.M, v.M, tc.name)
			assert.Equal(tc.expectedQueueValues, values, tc.name)
			assert.Equal(tc.expectedNode, tc.node, tc.name)
		})
	}
}

func addVisited(v *data.Visited, links ...string) *data.Visited {
	for _, s := range links {
		v.M[s] = true
	}
	return v
}
