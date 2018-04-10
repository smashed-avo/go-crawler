# Go Crawler

This simple web crawler is written in Go.

### Usage:

* Crawl URL
```
curl -X GET http://localhost:8000/crawl?url=https://medium.com/topic/technology
```

* Optional: Set crawling depth (Careful, it can slow down significantly)
```
curl -X GET http://localhost:8000/crawl?url=https://medium.com/topic/technology&depth=4
```

By default maximum crawling depth is set to 2, this means to get only first level children of seed URL

### Response

Following you can find an example response from the crawler in JSON format:

```
{
  "depth": 0,
  "title": "Not Found – Medium",
  "url": "https://medium.com/topic/technology~",
  "nodes": [
    {
      "depth": 1,
      "title": "Medium – Read, write and share stories that matter",
      "url": "https://medium.com/",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Medium – Read, write and share stories that matter",
      "url": "https://medium.com/",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Medium – Read, write and share stories that matter",
      "url": "https://medium.com/",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "",
      "url": "https://lukns.com/lost-in-antarctica-eef8fa970949?source=placement_card_footer_grid---------0-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "",
      "url": "https://lukns.com/lost-in-antarctica-eef8fa970949?source=placement_card_footer_grid---------0-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "",
      "url": "https://lukns.com/@ryanluikens",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "",
      "url": "https://lukns.com/@ryanluikens?source=placement_card_footer_grid---------0-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "To Anyone Who Has Lost Themselves: – Jamie Varon – Medium",
      "url": "https://medium.com/@jamievaron/to-anyone-who-has-lost-themselves-9c5e3049cb13?source=placement_card_footer_grid---------1-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "To Anyone Who Has Lost Themselves: – Jamie Varon – Medium",
      "url": "https://medium.com/@jamievaron/to-anyone-who-has-lost-themselves-9c5e3049cb13?source=placement_card_footer_grid---------1-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Jamie Varon – Medium",
      "url": "https://medium.com/@jamievaron",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Jamie Varon – Medium",
      "url": "https://medium.com/@jamievaron?source=placement_card_footer_grid---------1-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Three Things I Lost – Priya – Medium",
      "url": "https://medium.com/@priya_ebooks/three-things-i-lost-580108ca0a2c?source=placement_card_footer_grid---------2-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Three Things I Lost – Priya – Medium",
      "url": "https://medium.com/@priya_ebooks/three-things-i-lost-580108ca0a2c?source=placement_card_footer_grid---------2-31",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Priya – Medium",
      "url": "https://medium.com/@priya_ebooks",
      "nodes": []
    },
    {
      "depth": 1,
      "title": "Priya – Medium",
      "url": "https://medium.com/@priya_ebooks?source=placement_card_footer_grid---------2-31",
      "nodes": []
    }
  ]
}
```

### Requirements

* Golang 1.8+

* [dep](https://github.com/golang/dep)

### Libraries

The project relies mainly in Go standard libraries. There are a few dependencies added for the sake of simplicity and in areas where writing my own code does not add any value:

* [Goji.io](https://github.com/goji/goji) Minimalistic HTTP multiplexer to handle API calls.

* [PuerkitoBio/purell](github.com/PuerkitoBio/purell) URL sanitise - Go URL parse still accepts some links as valid and needed to sanitise further.

* [PuerkitoBio/goquery](github.com/PuerkitoBio/goquery) Website parsing library used to obtain website titles. Used as alternative to inspecting title element on tokenisation which was not giving the desired results. May revisit and remove this dependency because opening twice each site has a big impact on performance.

### Design considerations

* The project consists on the following packages:
```
.
├── cmd                          # cmd folder - contains the project executables
│   └── go-crawler               
│       └── main.go              # Main package and file - starts the server
└── lib                          # Application source code
    ├── crawler                  # Crawler package - Divided in two files, crawler.go is the parent process and worker.go crawls individual website
    │   └── crawler.go           # Process crawling seed URL and spins up the workers
    │   └── worker.go            # Creates website node, obtains title and spins up go routines to inspect web content and extract all links  
    ├── data                     # Data package
    │   └── data.go              # Contains Response struct used to store crawled info and unmarshal as JSON response to API call  
    ├── handler                  # Handler package
    │   └── handler.go           # Process seed URL and depth parameters and calls the crawling process  
    └── webpage                  # Webpage package
        └── links.go             # Loads the website, tokenises the DOM for the given URL and returns all links until end of document is reached
```

The main moving parts of the system are:

* crawler.go/Crawl - This is the parent process where most of the action happens:
  * Creates the parent node corresponding to the initial seed URL.
  * Fires an initial go Worker routine to process this first node.
  * Starts process loop where listen for new nodes added to the queue.
  * Fires a new worker for each node that has not been visited only if the maximum depth has not been reached.
  * Process only ends when all the workers have communicated they have finished.
  * Depth is maintained and passed to the workers so this info can be added to child nodes on creation.

* crawler.go/Worker - This is the child routine that process a single child URL:
  * Fires process to obtain all website links and then stays in a loop listening for:
    * Link - New link is created as child node and added to the array. Continues listening for new links.
    * Error - There was a problem opening the site and the process ends.
    * Done - Website parsing is complete and it communicates node array to parent process crawler.go/Crawl.
  * GetPageTitle() is called on node creation to open the page and obtain the web title.

* webpage.go/Links - This process is in charge of parsing the webpage and extract all links.
  * Opens an http client with a sensible timeout so if the site is unreachable, the process does not get stuck.
  * Starts a tokenisation process of the DOM to identify tags that contains an href link.
  * When an href link is found there are 2 levels of sanitisation happening:
    * Make sure the link starts with http*
    * Using a library make sure that it is not only a parseable URL but also a valid link, removing double slashes and fragments (hashlinks).

Synchronisation of the different routines at the two levels: Crawl<->Worker and Worker<->Links occurs via channels. These channel reads are blocking and sync the execution of the threads. This follows the paradigm in Go of blocking by communicating instead of by shared memory.

### Getting Started

Install the dependent libraries using dep

```
dep ensure
```

### Running the service

Then start the server by running:
```
go run cmd/go-crawler/main.go
```

The application runs by default in http://localhost:8000, on a later stage configuration can be added to modify this based on environment variables.

### Testing

Unit tests are missing and are planned to be added soon. Packages need interfaces so unit tests can use mocking of dependencies.

More TBD
