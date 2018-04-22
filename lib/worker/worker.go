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
	Collector  Collectorer
	ChLinks    chan string
	ChErrors   chan error
	ChFinished chan bool
}

// NewWorker factory method to inject collector instance
func NewWorker(c Collectorer) *Worker {
	return &Worker{Collector: c}
}

// Do gets all links for a website and stores them in the node
func (w *Worker) Do(node *data.Response, depth int, chQueue chan []*data.Response, visited *data.Visited) {
	// Channels
	w.ChLinks = make(chan string)
	w.ChFinished = make(chan bool)
	w.ChErrors = make(chan error)
	defer close(w.ChLinks)
	defer close(w.ChFinished)
	defer close(w.ChErrors)

	go w.Collector.Collect(node.URL, w.ChLinks, w.ChFinished, w.ChErrors)

	// Subscribe to channels
	for {
		select {
		case link := <-w.ChLinks:
			// check if link ir parseable
			u, err := url.Parse(link)
			if err != nil {
				println(err.Error())
				continue
			}
			visited.RLock()
			// If node already visited, do not register
			if visited.M[link] {
				continue
			}
			visited.M[link] = true
			visited.RUnlock()
			subNode := data.Response{
				Depth: depth,
				Title: w.GetPageTitle(u.String()),
				URL:   u.String(),
				Nodes: make([]*data.Response, 0),
			}
			node.Nodes = append(node.Nodes, &subNode)
			break
		case <-w.ChFinished:
			chQueue <- node.Nodes
			return
		case <-w.ChErrors:
			// Failed to return 200 OK for this link
			chQueue <- node.Nodes
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
