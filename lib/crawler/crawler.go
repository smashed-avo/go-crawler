package crawler

import (
	"fmt"
	"net/url"

	"github.com/smashed-avo/go-crawler/lib/data"
)

// FetchService service to initiate URL crawling
//type FetchService struct {
//}

//type Crawler interface {
//	Crawl(url string, ch chan string, chFinished chan bool)
//}

// Fetch Initiates crawl process given a set of seed URLs
func Crawl(seedUrl *url.URL, maxDepth int) (response *data.Response, err error) {

	//setup channel for inputs to be processed
	pending := make(chan *data.Response, 0)
	visited := make(map[string]bool)
	finished := make(chan *data.Response, 0)
	defer close(pending)
	defer close(finished)

	// Start background crawling processor
	go process(maxDepth, pending, visited, finished)

	// Initialise channel & start crawling
	site := data.Response{
		Depth: 0,
		Title: seedUrl.String(),
		Nodes: make([]*data.Response, 0),
		URL:   seedUrl.String(),
	}
	pending <- &site

	// Wait to finish
	res := <-finished

	// TODO Build response

	fmt.Println("\n--- Completed ---")
	for link := range visited {
		fmt.Println("Visited: " + link)
	}

	return res, nil
}

func process(maxDepth int, pending chan *data.Response, visited map[string]bool, finished chan *data.Response) {
	depth := 1
	var processed *data.Response
	for depth > 0 {
		next := <-pending
		depth--
		processed = next

		// if were too deep, skip it
		if next.Depth >= maxDepth {
			continue
		}

		//loop over all urls to visit from that page
		for _, node := range next.Nodes {

			//check we haven't visited them before
			if visited[node.URL] {
				continue
			}

			// update control structures
			depth++
			visited[node.URL] = true

			go Worker(node, depth, pending)
		}
	}
	finished <- processed
}
