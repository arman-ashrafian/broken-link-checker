package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

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

	ch := make(chan string)
	for _, l := range links {
		if l[0] == '/' {
			l = baseURL + l
		} else if l[0:4] != "http" {
			l = baseURL + "/" + l
		}
		go makeRequest(l, ch)
	}

	for x := 0; x < 1000; x++ {
		v, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("%d --", x)
		fmt.Println(v)
	}
	fmt.Println("done")
}

func makeRequest(link string, ch chan<- string) {
	resp, err := http.Get(link)
	if err != nil || resp.StatusCode != 200 {
		ch <- fmt.Sprintf("%s  --  ERROR", link)
		return
	}
	ch <- fmt.Sprintf("%s  --  %s", link, resp.Status)
}
