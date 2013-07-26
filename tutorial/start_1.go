package main

import (
	"flag"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"log"
	"os"
)

var (
	host *string = flag.String("host", "localhost", "Elasticsearch Host")
)

func main() {
	core.DebugRequests = true
	log.SetFlags(log.LstdFlags)
	flag.Parse()

	fmt.Println("host = ", *host)
	// Set the Elasticsearch Host to Connect to
	api.Domain = *host

	// Index a document
	_, err := core.Index(false, "testindex", "user", "docid_1", `{"name":"bob"}`)
	exitIfErr(err)

	// Index a doc using a map of values
	_, err = core.Index(false, "testindex", "user", "docid_2", map[string]string{"name": "venkatesh"})
	exitIfErr(err)

	// Index a doc using Structs
	_, err = core.Index(false, "testindex", "user", "docid_3", MyUser{"wanda", 22})
	exitIfErr(err)

	// Search Using Raw json String
	searchJson := `{
	    "query" : {
	        "term" : { "name" : "wanda" }
	    }
	}`
	out, err := core.SearchRequest(true, "testindex", "user", searchJson, "")
	if len(out.Hits.Hits) == 1 {
		fmt.Println(string(out.Hits.Hits[0].Source))
	}
	exitIfErr(err)

}
func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

type MyUser struct {
	Name string
	Age  int
}
