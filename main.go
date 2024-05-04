package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Story struct {
	By          string `json:"by"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Kids        []int  `json:"kids"`
	Descendants int    `json:"descendants"`
	Id          int    `json:"id"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
}

func main() {
	endpoint := "https://hacker-news.firebaseio.com/v0/"
	topstories := getTopStories(endpoint)

	for _, v := range topstories {
		story := getStory(endpoint, v)
		fmt.Println(story.Title)
	}
}

func getTopStories(endpoint string) []int {
	url := endpoint + "topstories.json"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	// We Read the response body on the line below
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	topstories := make([]int, 0, 500)
	json.Unmarshal(body, &topstories)

	return topstories
}

func getStory(endpoint string, id int) Story {
	url := endpoint + "item/" + fmt.Sprint(id) + ".json"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var story Story
	json.Unmarshal(body, &story)
	return story
}
