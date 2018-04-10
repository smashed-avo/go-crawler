package crawler

import (
	"net/url"

	"github.com/smashed-avo/go-crawler/lib/data"
)

// Crawl Initiates crawl process given a set of seed URLs
func Crawl(seedURL *url.URL, maxDepth int) (response *data.Response) {
	//setup channels to process nodes recursively
	chQueue := make(chan []*data.Response)
	defer close(chQueue)

	// Maintain visited URL to detect loops
	visited := make(map[string]bool)

	// add first parent node to queue
	parent := data.Response{
		Depth: 0,
		Title: GetPageTitle(seedURL.String()),
		Nodes: make([]*data.Response, 0),
		URL:   seedURL.String(),
	}
	depth := 1
	workers := 1
	go Worker(&parent, depth, chQueue)

	for workers > 0 {
		nodes := <-chQueue
		workers--
		if len(nodes) > 0 {
			depth = nodes[0].Depth + 1
			if depth < maxDepth {
				for _, node := range nodes {
					if visited[node.URL] {
						continue
					}
					visited[node.URL] = true
					workers++
					go Worker(node, depth, chQueue)
				}
			}
		}
	}
	return &parent
}
