package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

// info to output from HTTP GET request
type responseInfo struct {
	link   string
	status int
}

const baseURL = "http://computingpaths.ucsd.edu"

var links []string
var repeatedLinks int // count repeats

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

func main() {
	cache := make(map[string]bool)
	repeatedLinks = 0

	fmt.Println("Checking Links")
	fmt.Println("--------------")

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

		// skip to next iteration if link has already been checked
		_, ok := cache[l]
		if ok {
			repeatedLinks += 1
			continue
		}

		go makeRequest(l, ch)
		cache[l] = true
	}

	// display failed requests and url to file
	f, _ := os.Create("deadlinks.txt")
	defer f.Close()
	for x := 0; x < len(links)-repeatedLinks; x++ {
		v := <-ch
		if v.status == 200 {
			continue
		}
		fmt.Printf("|%-3d| %-6s\n", v.status, v.link)
		fmt.Fprintf(f, "%s\n", v.link)
	}
	fmt.Println(" DONE")
}
