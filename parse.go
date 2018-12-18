package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type DB struct {
	Pages   []Page  `json:"pages"`
	Stories []Story `json:"stories"`
	Majors  []Major `json:"majors"`
}

type Page struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

type Story struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

type Major struct {
	Name     string          `json:"name"`
	MoreInfo []MajorMoreInfo `json:"moreInfo"`
	Image    string          `json:"image"`
}

type MajorMoreInfo struct {
	Name string
	Link string
}

func parse() DB {
	jsonfile, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	jsonfile_bytes, _ := ioutil.ReadAll(jsonfile)

	var db DB
	json.Unmarshal(jsonfile_bytes, &db)

	return db
}

func getLinks() []string {
	db := parse()
	var links []string

	for _, p := range db.Pages {
		links = append(links, p.Link)
	}
	for _, s := range db.Stories {
		links = append(links, s.Image)
		links = append(links, s.Link)
	}
	for _, m := range db.Majors {
		links = append(links, m.Image)
		for _, info := range m.MoreInfo {
			links = append(links, info.Link)
		}
	}
	return links
}

func main() {
	links := getLinks()
	for _, l := range links {
		fmt.Println(l)
	}
	fmt.Println(len(links))
}
