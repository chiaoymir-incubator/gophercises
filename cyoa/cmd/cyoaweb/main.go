package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/woodcutter-eric/gophercises/cyoa/cyoa"
)

func main() {
	port := flag.Int("port", 3001, "the server port")
	filename := flag.String("file", "gopher.json", "the story file")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	story, err := cyoa.ParseJSONStory(f)
	if err != nil {
		log.Fatal(err)
		return
	}

	// %v: print map content
	// %+v: print map key and content
	// fmt.Printf("%+v\n", story)

	h := cyoa.NewHandler(story)
	addr := fmt.Sprintf(":%d", *port)

	fmt.Printf("Starting the server on %d\n", *port)
	log.Fatal(http.ListenAndServe(addr, h))

}
