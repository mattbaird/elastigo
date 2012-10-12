package main

import (
	"github.com/mattbaird/elastigo/core"
	"log"
)

// for testing
func main() {
	//	core.RunSearch(true, "questionscraper", "user:kimchy")
	log.Println("Index:")
	response, _ := core.Index(true, "twitter", "tweet", "1", NewTweet("kimchy", "Search is cool"))
	log.Printf("Get: %s", response.Exists)
	response, _ = core.Get(true, "twitter", "tweet", "1")
	log.Println("Exists:")
	response, _ = core.Exists(true, "twitter", "tweet", "1")
	log.Println("Delete:")
	response, _ = core.Delete(true, "twitter", "tweet", "1", -1, "")
	response, _ = core.Get(true, "twitter", "tweet", "1")
	log.Printf("Get: %s", response.Exists)
}
