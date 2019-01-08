package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

var db = parse()

func parse() DB {
	jsonfile, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	jsonfileBytes, _ := ioutil.ReadAll(jsonfile)

	var db DB
	json.Unmarshal(jsonfileBytes, &db)

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
		links = append(links, b.Image)
		links = append(links, b.Link)
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
		checkLinkAndUpdate(link, brokenLink, p, &db.Pages[i].Link)
	}

	// check if brokenLink is in Stories
	for i, s := range db.Stories {
		link := normalizeLink(s.Link)
		image := normalizeLink(s.Image)
		checkLinkAndUpdate(link, brokenLink, s, &db.Stories[i].Link)
		checkLinkAndUpdate(image, brokenLink, s, &db.Stories[i].Image)
	}

	// check if brokenLink is in Majors
	for i, m := range db.Majors {
		image := normalizeLink(m.Image)
		checkLinkAndUpdate(image, brokenLink, m, &db.Majors[i].Image)

		for j, info := range m.MoreInfo {
			link := normalizeLink(info.Link)
			checkLinkAndUpdate(link, brokenLink, m, &db.Majors[i].MoreInfo[j].Link)
		}
	}

	// check if brokenLink is in Departments
	for i, d := range db.Departments {
		link := normalizeLink(d.(map[string]interface{})["link"].(string))
		if link == brokenLink {
			newLink := promptUserForUpdate(d, brokenLink)
			if newLink != "" {
				db.Departments[i].(map[string]interface{})["link"] = newLink
			}
		}

	}

}

func checkLinkAndUpdate(link, brokenLink string, data interface{}, dbLink *string) {
	newLink := ""
	if link == brokenLink {
		newLink = promptUserForUpdate(data, brokenLink)
		if newLink != "" {
			*dbLink = newLink
		}
	}
}

// prompt user to update broken link
// returns the updated link or empty string if user wants to skip
func promptUserForUpdate(data interface{}, brokenLink string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	prettyPrint(data)
	fmt.Println()

	fmt.Printf("Broken Link: %s\n", brokenLink)

	fmt.Print("update link: ")
	resp, _ := reader.ReadString('\n')
	resp = strings.TrimSuffix(resp, "\n") // remove '\n'

	if resp == "x" || resp == "X" {
		return ""
	}

	return resp
}

func saveDBToFile(filename string, db DB) {
	b, _ := json.MarshalIndent(db, "", "  ")
	f, _ := os.Create(filename)
	f.Write(b)
	f.Close()
}
