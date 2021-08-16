package utils

import (
	"encoding/json"
	"io"
	"log"
)

type Story map[string]Chapter

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

func ParseJSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	story := Story{}

	err := d.Decode(&story)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return story, nil
}
