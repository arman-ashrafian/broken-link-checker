package main

import (
	"fmt"
	"net/http"
	"time"
)

// info to output from HTTP GET request
type responseInfo struct {
	link   string
	status int
}

const baseURL = "http://computingpaths.ucsd.edu"

var brokenLinks []string
var repeatedLinks int // count repeats

func makeRequest(link string, ch chan<- *responseInfo) {
	timeout := time.Duration(time.Second * 10)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(link)
	info := new(responseInfo)
	info.link = link

	if err != nil {
		info.status = -1
		ch <- info
		return
	} else if resp.StatusCode != 200 {
		info.status = resp.StatusCode
		ch <- info
		return
	}
	info.status = 200
	ch <- info
}

func normalizeLink(l string) string {
	if l == "" {
		return ""
	}

	if l[0] == '/' {
		l = baseURL + l
	} else if l[0:4] != "http" {
		l = baseURL + "/" + l
	}
	return l
}

func checkLinks() {
	cache := make(map[string]bool)
	repeatedLinks = 0

	fmt.Println("Checking Links")
	fmt.Println("--------------")

	links := getLinks()

	// sync responses from all GET requests
	ch := make(chan *responseInfo)

	// make GET requests
	for _, l := range links {
		if l == "" {
			repeatedLinks += 1
			continue
		}

		// add basepath to link if not there
		l = normalizeLink(l)

		// skip to next iteration if link has already been checked
		_, ok := cache[l]
		if ok {
			repeatedLinks += 1
			continue
		}

		go makeRequest(l, ch)
		cache[l] = true
	}

	// display failed requests
	for x := 0; x < len(links)-repeatedLinks; x++ {
		v := <-ch
		if v.status != 200 {
			fmt.Printf("|%-3d| %-6s\n", v.status, v.link)
			brokenLinks = append(brokenLinks, v.link)
		}
	}
	fmt.Println("...DONE")
}

func main() {
	updateLink("http://computingpaths.ucsd.edu/stories.html")
	updateLink("https://www.youtube.com/watch?v=AAcQylywrYE")
	updateLink("http://www.be.ucsd.edu/bioinformatics")
	saveDBToFile("database_out.json", db)
}
