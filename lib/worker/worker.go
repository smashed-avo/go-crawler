package worker

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/smashed-avo/go-crawler/lib/data"
)

// Collectorer interface to collector function
type Collectorer interface {
	Collect(url string, chLinks chan string, chFinished chan bool, chErrors chan error)
}

// Worker responsible for a single URL to retrieve all its linkr and store them as linked nodes
type Worker struct {
	Collector Collectorer
}

// NewWorker factory method to inject ollector instance
func NewWorker(c Collectorer) *Worker {
	return &Worker{Collector: c}
}

// Do gets all links for a website and stores them in the node
func (w *Worker) Do(node *data.Response, depth int, chQueue chan []*data.Response, visited *data.Visited) {
	// Channels
	chLinks := make(chan string)
	chFinished := make(chan bool)
	chErrors := make(chan error)
	defer close(chLinks)
	defer close(chFinished)
	defer close(chErrors)

	// fetcher := getFetcher()
	go w.Collector.Collect(node.URL, chLinks, chFinished, chErrors)

	// Subscribe to channels
	for {
		select {
		case link := <-chLinks:
			visited.RLock()
			// If node already visited, do not register
			if visited.M[node.URL] {
				continue
			}
			visited.RUnlock()
			u, err := url.Parse(link)
			if err != nil {
				println(err.Error())
			}
			subNode := data.Response{
				Depth: depth,
				Title: w.GetPageTitle(u.String()),
				URL:   u.String(),
				Nodes: make([]*data.Response, 0),
			}
			node.Nodes = append(node.Nodes, &subNode)
			break
		case <-chFinished:
			chQueue <- node.Nodes
			return
		case <-chErrors:
			// Failed to return 200 OK for this link
			return
		}
	}
}

// GetPageTitle gets title for a page
func (w *Worker) GetPageTitle(u string) string {
	doc, err := goquery.NewDocument(u)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find("title").Text())
}
