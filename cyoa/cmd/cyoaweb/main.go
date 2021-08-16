package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/woodcutter-eric/gophercises/cyoa/utils"
)

func main() {
	filename := flag.String("file", "gopher.json", "the story file")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	story, err := utils.ParseJSONStory(f)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%+v\n", story)
}
