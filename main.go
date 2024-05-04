package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Item struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Children  []int  `json:"children"`
	Id        int    `json:"id"`
	Points    int    `json:"points"`
}

type SearchResult struct {
	Hits []Item `json:"hits"`
}

func main() {
	endpoint := "http://hn.algolia.com/api/v1/"
	stories := getFrontPage(endpoint)
	for _, v := range stories {
		fmt.Println(v.Title)
	}
}

func getFrontPage(endpoint string) []Item {
	url := endpoint + "search?tags=front_page"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	// We Read the response body on the line below
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var searchResult SearchResult
	err2 := json.Unmarshal(body, &searchResult)
	if err2 != nil {
		fmt.Println(err2)
	}

	return searchResult.Hits
}

func getItem(endpoint string, id int) Item {
	url := endpoint + "items/" + fmt.Sprint(id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var item Item
	json.Unmarshal(body, &item)
	return item
}
