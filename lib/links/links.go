package links

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
	"golang.org/x/net/html"
)

// Collector processes a webpage and collect all links
type Collector struct {
	client *http.Client
}

// New returns a pointer to a new collector
func NewCollector(client *http.Client) *Collector {
	return &Collector{client: client}
}

// Links extract title and all links from a given URL
func (c *Collector) Collect(url string, chLinks chan string, chFinished chan bool, chErrors chan error) {

	resp, err := c.client.Get(url)

	if err != nil {
		chErrors <- err
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// End of the document
			chFinished <- true
			return
		case html.StartTagToken, html.EndTagToken:
			t := z.Token()
			if "a" == t.Data {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						// Make sure the url begins with http
						if strings.Index(attr.Val, "http") == 0 {
							if val, err := normalizeURL(attr.Val); err == nil {
								chLinks <- val
							}
						}
					}
				}
			}
		}
	}
}

func normalizeURL(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return purell.NormalizeURL(parsedURL, purell.FlagsSafe|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveFragment), nil
}