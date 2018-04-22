package main

import (
	"fmt"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/smashed-avo/go-crawler/lib/crawler"
	"github.com/smashed-avo/go-crawler/lib/handler"
	"github.com/smashed-avo/go-crawler/lib/links"
	"github.com/smashed-avo/go-crawler/lib/worker"
)

const port = "8000"

// main sets the router and starts the serves
func main() {
	h := getHandler()
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/crawl"), h.HandleCrawl)
	fmt.Printf("Server started. Listening on port %s.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getHandler() *handler.Handler {
	client := &links.HTTPClient{}
	l := links.NewCollector(client)
	w := worker.NewWorker(l)
	c := crawler.NewCrawler(w)

	return handler.NewHandler(c)
}
