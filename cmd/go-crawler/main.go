package main

import (
	"fmt"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/smashed-avo/go-crawler/lib/handler"
)

const port = "8000"

// main sets the router and starts the serves
func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/crawl"), handler.HandleCrawl)
	fmt.Printf("Server started. Listening on port %s.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
