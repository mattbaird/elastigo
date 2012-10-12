package main

import (
	"github.com/mattbaird/elastigo/core"
)

// for testing
func main() {
	//	core.RunSearch(true, "questionscraper", "user:kimchy")
	core.Index(true, "twitter", "tweet", "1", NewTweet("kimchy", "Search is cool"))
	core.Get(true, "twitter", "tweet", "1")
	core.Exists(true, "twitter", "tweet", "1")
}
