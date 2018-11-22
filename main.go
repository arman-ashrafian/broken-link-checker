package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

type responseObj struct {
	link   string
	status int
}

func main() {
	const baseURL = "http://computingpaths.ucsd.edu"

	var links []string
	fmt.Println("Checking Broken Links")
	fmt.Println("---------------------")

	file, err := os.Open("links.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}

	ch := make(chan *responseObj)
	for _, l := range links {
		if l[0] == '/' {
			l = baseURL + l
		} else if l[0:4] != "http" {
			l = baseURL + "/" + l
		}
		go makeRequest(l, ch)
	}

	for x := 0; x < len(links); x++ {
		v := <-ch
		if v.status != 200 {
			fmt.Printf("|%s| %-6s\n", "ERROR", v.link)
		}
	}
	fmt.Println("done")
}

func makeRequest(link string, ch chan<- *responseObj) {
	resp, err := http.Get(link)
	respObj := new(responseObj)
	respObj.link = link

	if err != nil || resp.StatusCode != 200 {
		respObj.status = -1
		ch <- respObj
		return
	}
	respObj.status = 200
	ch <- respObj
}
