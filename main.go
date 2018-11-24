package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

// info to output from HTTP GET request
type responseInfo struct {
	link   string
	status int
}

const baseURL = "http://computingpaths.ucsd.edu"

var links []string

func readLinksFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
}

func makeRequest(link string, ch chan<- *responseInfo) {
	resp, err := http.Get(link)
	info := new(responseInfo)
	info.link = link

	if err != nil || resp.StatusCode != 200 {
		info.status = -1
		ch <- info
		return
	}
	info.status = 200
	ch <- info
}

func main() {
	fmt.Println("Checking Broken Links")
	fmt.Println("---------------------")

	readLinksFromFile("links.txt")

	// sync responses from all GET requests
	ch := make(chan *responseInfo)

	// make GET requests
	for _, l := range links {
		if l[0] == '/' {
			l = baseURL + l
		} else if l[0:4] != "http" {
			l = baseURL + "/" + l
		}
		go makeRequest(l, ch)
	}

	// display failed requests
	for x := 0; x < len(links); x++ {
		v := <-ch
		if v.status != 200 {
			fmt.Printf("|%s| %-6s\n", "ERROR", v.link)
		}
	}
	fmt.Println(" DONE")
}
