package crawler

import (
	"net/url"
	"strings"

	"github.com/smashed-avo/go-crawler/lib/data"
	"github.com/smashed-avo/go-crawler/lib/webpage"

	"github.com/PuerkitoBio/goquery"
)

// Worker gets all links for a website and stores it in the node
func Worker(node *data.Response, depth int, chQueue chan []*data.Response) {
	// Channels
	chLinks := make(chan string)
	chFinished := make(chan bool)
	chErrors := make(chan error)
	defer close(chLinks)
	defer close(chFinished)
	defer close(chErrors)

	go webpage.Links(node.URL, chLinks, chFinished, chErrors)

	// Subscribe to channels
	for {
		select {
		case link := <-chLinks:
			u, err := url.Parse(link)
			if err != nil {
				println(err.Error())
			}
			subNode := data.Response{
				Depth: depth,
				Title: GetPageTitle(u.String()),
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
func GetPageTitle(u string) string {
	doc, err := goquery.NewDocument(u)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find("title").Text())
}
