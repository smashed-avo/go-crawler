package handler

import (
	"fmt"
	"net/http"

	"encoding/json"
	"net/url"

	"github.com/smashed-avo/go-crawler/lib/crawler"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

// Fetcher returns the response containing all crawled for a URL
//type Fetcher interface {
//	Fetch(url *url.URL) (response *data.Response, err error)
//}
//
//var urlFetcher Fetcher

// HandleCrawl handles the crawl api request
func HandleCrawl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)

	// Get URL parameter and validate/sanitise
	urlParam := r.URL.Query().Get("url")
	fmt.Println("Seed URL: " + urlParam)
	u, err := url.ParseRequestURI(urlParam)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO accept depth param

	//Start crawling process
	resp, err := crawler.Crawl(u, 10)
	if err != nil {
		println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
