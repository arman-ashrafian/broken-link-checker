package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type DB struct {
	Pages          []Page      `json:"pages"`
	Stories        []Story     `json:"stories"`
	Majors         []Major     `json:"majors"`
	Departments    Departments `json:"departments"`
	Resources      []Resource  `json:"resources"`
	ResourceBanner []ResBanner `json:"resourceBanner"`
	Projects       []Project   `json:"projects"`
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
	Name string `json:"name"`
	Link string `json:"link"`
}

type Departments map[string]interface{}
type Department struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type Resource struct {
	Name     string `json:"name"`
	MapImage string `json:"mapImage,omitempty"`
	MapLink  string `json:"mapLink,omitempty"`
	Link     string `json:"link"`
}

type ResBanner struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

type Project struct {
	Name   string   `json:"name"`
	Images []string `json:"images"`
	Videos []string `json:"videos"`
	Link   string   `json:"link"`
}

var db DB = parse()

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
	for _, v := range db.Departments {
		links = append(links, v.(map[string]interface{})["link"].(string))
	}
	for _, r := range db.Resources {
		links = append(links, r.Link)
		links = append(links, r.MapLink)
		links = append(links, r.MapImage)
	}
	for _, b := range db.ResourceBanner {
		links = append(links, r.Image)
		links = append(links, r.Link)
	}

	return links
}

func prettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func updateLink(brokenLink string) {
	// check if brokenLink is in Pages
	for i, p := range db.Pages {
		link := normalizeLink(p.Link)
		if link == brokenLink {
			newLink := promptUserForUpdate(p, brokenLink)
			db.Pages[i].Link = newLink
		}

	}

	// check if brokenLink is in Stories
	for _, s := range db.Stories {
		link := normalizeLink(s.Link)
		image := normalizeLink(s.Image)
		if link == brokenLink {
			newLink := promptUserForUpdate(s, brokenLink)
			s.Link = newLink
		}
		if image == brokenLink {
			newLink := promptUserForUpdate(s, brokenLink)
			s.Image = newLink
		}
	}

	// check if brokenLink is in Majors
	for _, m := range db.Majors {
		image := normalizeLink(m.Image)
		if image == brokenLink {
			newLink := promptUserForUpdate(m, brokenLink)
			m.Image = newLink
		}
		for _, info := range m.MoreInfo {
			link := normalizeLink(info.Link)
			if link == brokenLink {
				newLink := promptUserForUpdate(m, brokenLink)
				info.Link = newLink
			}
		}
	}
}

func promptUserForUpdate(data interface{}, brokenLink string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	prettyPrint(data)
	fmt.Println()

	fmt.Printf("Broken Link: %s\n", brokenLink)

	fmt.Print("update link: ")
	resp, _ := reader.ReadString('\n')
	resp = resp[:len(resp)-2] // remove '\n'

	if resp == "x\n" || resp == "X\n" {
		return ""
	}

	return resp
}
