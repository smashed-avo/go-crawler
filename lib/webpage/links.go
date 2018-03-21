package webpage

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"time"
)

// Links extract all links from a given URL
func Links(url string, links chan string, chFinished chan bool) {

	// Set timeout to 15s
	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := c.Get(url)

	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		chFinished <- true
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
							links <- attr.Val
						}
					}
				}
			}
		}
	}
}
