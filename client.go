package main

import (
	"github.com/mattbaird/elastigo/core"
)

// for testing
func main() {
	//	core.RunSearch(true, "questionscraper", "user:kimchy")
	core.RunGet(true, "questionscraper", "question_holdingpen", "50760a78d23ea1e0f51d52de")
}
