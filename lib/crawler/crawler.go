package crawler

import (
	"net/url"

	"github.com/smashed-avo/go-crawler/lib/data"
)

// Workerer is an interface to the worker function
type Workerer interface {
	Do(node *data.Response, depth int, chQueue chan []*data.Response)
	GetPageTitle(u string) string
}

// Crawler receiver for crawl function
type Crawler struct {
	Worker  Workerer
	ChQueue chan []*data.Response
}

// NewCrawler factory method to inject worker instance
func NewCrawler(w Workerer) *Crawler {
	return &Crawler{Worker: w}
}

// Crawl Initiates crawl process given an initial seed URLs
func (c *Crawler) Crawl(seedURL *url.URL, maxDepth int) *data.Response {
	//setup channels to process nodes recursively
	c.ChQueue = make(chan []*data.Response)
	defer close(c.ChQueue)

	// Maintain visited URL to detect loops
	visited := make(map[string]bool)

	// add first parent node to queue
	parent := data.Response{
		Depth: 0,
		Title: c.Worker.GetPageTitle(seedURL.String()),
		Nodes: make([]*data.Response, 0),
		URL:   seedURL.String(),
	}
	depth := 1
	workers := 1
	go c.Worker.Do(&parent, depth, c.ChQueue)

	for workers > 0 {
		nodes := <-c.ChQueue
		workers--
		if len(nodes) > 0 {
			depth = nodes[0].Depth + 1
			// println(depth)
			// println(" , ")
			// println(maxDepth)
			if depth < maxDepth {
				for _, node := range nodes {
					if visited[node.URL] {
						continue
					}
					visited[node.URL] = true
					workers++
					go c.Worker.Do(node, depth, c.ChQueue)
				}
			}
		}
	}
	return &parent
}
