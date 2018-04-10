package webpage

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/purell"
	"golang.org/x/net/html"
)

// Links extract title and all links from a given URL
func Links(url string, chLinks chan string, chFinished chan bool, chErrors chan error) {
	// Set timeout to 15s
	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := c.Get(url)

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
