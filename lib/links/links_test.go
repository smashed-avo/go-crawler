package links_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/smashed-avo/go-crawler/lib/links"
	"github.com/stretchr/testify/assert"
)

var (
	threeLinksHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>title</title>
    <link rel="stylesheet" href="style.css">
    <script src="script.js"></script>
  </head>
  <body>
		 <a href="https://www.linkedsite1.com">Visit me!</a>
		 <a href="https://www.linkedsite2.com">Visit me!</a>
		 <a href="https://www.linkedsite3.com">Visit me!</a>
  </body>
</html>`
	nonParseableLinkHTML = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>title</title>
		<link rel="stylesheet" href="style.css">
		<script src="script.js"></script>
	</head>
	<body>
		<a href="http://a b.com/">Visit me!</a>
	</body>
</html>`
	fragmentLinkHTML = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>title</title>
		<link rel="stylesheet" href="style.css">
		<script src="script.js"></script>
	</head>
	<body>
		<a href="https://www.linkedsite1.com">Visit me!</a>
 		<a href="#div_id">jump link</a>
 		<div id="div_id">jump here</div>
	</body>
</html>`
)

const (
	_ mockStateClient = iota
	success
	nonParseableLink
	sanitiseFragment
	errorClient
)

type mockStateClient int

type MockClient struct {
	State mockStateClient
}

func (c *MockClient) Get(url string) (io.Reader, error) {
	switch c.State {
	case success:
		return strings.NewReader(threeLinksHTML), nil
	case nonParseableLink:
		return strings.NewReader(nonParseableLinkHTML), nil
	case sanitiseFragment:
		return strings.NewReader(fragmentLinkHTML), nil
	case errorClient:
		return nil, errors.New(`couldn't fetch website`)
	default:
		panic(fmt.Sprintf("Invalid mockStateClient: %v", c.State))
	}
}

func TestDo(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		name          string
		state         mockStateClient
		expectedLinks []string
		expectedError error
	}{
		{
			name:          "Success - Collect three links",
			state:         success,
			expectedLinks: []string{`https://www.linkedsite1.com`, `https://www.linkedsite2.com`, `https://www.linkedsite3.com`},
			expectedError: nil,
		},
		{
			name:          "Success - Link non parseable ommited",
			state:         nonParseableLink,
			expectedLinks: []string{},
			expectedError: nil,
		},
		{
			name:          "Success - Hash link sanitised",
			state:         sanitiseFragment,
			expectedLinks: []string{`https://www.linkedsite1.com`},
			expectedError: nil,
		},
		{
			name:          "Error - client fetch failed",
			state:         errorClient,
			expectedLinks: []string{},
			expectedError: errors.New(`couldn't fetch website`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := MockClient{State: tc.state}
			c := links.NewCollector(&m)

			chLinks := make(chan string)
			chErrors := make(chan error)
			chFinished := make(chan bool)

			go c.Collect(`www.google.com`, chLinks, chFinished, chErrors)

			links := make([]string, 0)
			var err error
		loop:
			for {
				select {
				case link := <-chLinks:
					links = append(links, link)
				case <-chFinished:
					break loop
				case err = <-chErrors:
					break loop
				}
			}

			close(chLinks)
			close(chErrors)
			close(chFinished)

			assert.Equal(tc.expectedLinks, links, tc.name)
			assert.Equal(tc.expectedError, err, tc.name)
		})
	}
}
