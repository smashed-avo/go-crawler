package crawler

import (
	"fmt"
	"net/url"

	"github.com/smashed-avo/go-crawler/lib/data"
	"github.com/smashed-avo/go-crawler/lib/webpage"
)

// Worker gets all links for a website and stores it in the node
func Worker(parent *data.Response, depth int, pending chan *data.Response) {
	// Channels
	links := make(chan string)
	finished := make(chan bool)

	go webpage.Links(parent.URL, links, finished)

	// Subscribe to both channels
	select {
	case u := <-links:
		fmt.Println("Found url: " + u)
		site, err := url.Parse(u)
		if err != nil {
			println(err.Error())
		}
		pending <- &data.Response{
			Depth: depth,
			Title: u,
			URL:   site.String(),
			Nodes: make([]*data.Response, 0),
		}
	case <-finished:
		return
	}
}
