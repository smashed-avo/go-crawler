package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/smashed-avo/go-crawler/lib/crawler"
)

// HandleCrawl handles the crawl api request
func HandleCrawl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get URL parameter and validate/sanitise
	urlParam := r.URL.Query().Get("url")
	u, err := url.ParseRequestURI(urlParam)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Default param value
	maxDepth := 2
	maxDepthParam := r.URL.Query().Get("depth")
	if maxDepthParam != "" {
		maxDepth, err = strconv.Atoi(maxDepthParam)
		if err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//Start crawling process
	res := crawler.Crawl(u, maxDepth)

	json.NewEncoder(w).Encode(res)
}
