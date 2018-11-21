package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
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
		if l[0:3] == "htt" {
			go makeRequest(l, ch)
		}
	}

	for _ = range links {
		fmt.Println(<-ch)
	}
}

func makeRequest(link string, ch chan<- string) {
	resp, _ := http.Get(link)
	ch <- fmt.Sprintf("%s  --  %s", link, resp.Status)
}
