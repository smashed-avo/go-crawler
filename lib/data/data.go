package data

import "net/url"

// Response model the response to API call
type Response struct {
	Depth int		  `json:"depth" description:"Depth of URL from seed website"`
	Nodes []*Response `json:"nodes" description:"Children of a site fetched by the crawler"`
	Title string      `json:"title" description:"Title of a site fetched by the crawler"`
	URL   string      `json:"url" description:"URL of a site fetched by the crawler"`
}

// Sites holds the visited urls and their relationships parent/children
type Sites struct {
	Depth  int
	Parent *url.URL
	URLs   []*url.URL
}

// Error encapsulates error message in body response
type Error struct {
	Message string `json:"message" description:"Error message."`
}
