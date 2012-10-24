package main

import (
	"encoding/json"
	"github.com/mattbaird/elastigo/cluster"
	"github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/indices"
	"log"
)

// for testing
func main() {
	response, _ := core.Index(true, "twitter", "tweet", "1", NewTweet("kimchy", "Search is cool"))
	indices.Flush()
	log.Printf("Index OK: %b", response.Ok)
	searchresponse, err := core.Search(true, "twitter", "tweet", "{\"query\" : {\"term\" : { \"user\" : \"kimchy\" }}}")
	if err != nil {
		log.Println("error during search:" + err.Error())
		log.Fatal(err)
	}
	// try marshalling to tweet type
	var t Tweet
	json.Unmarshal(searchresponse.Hits.Hits[0].Source, t)
	log.Printf("Search Found: %s", t)
	response, _ = core.Get(true, "twitter", "tweet", "1")
	log.Printf("Get: %s", response.Exists)
	response, _ = core.Exists(true, "twitter", "tweet", "1")
	log.Printf("Exists: %s", response.Exists)
	indices.Flush()
	countResponse, _ := core.Count(true, "twitter", "tweet")
	log.Printf("Count: %s", countResponse.Count)
	response, _ = core.Delete(true, "twitter", "tweet", "1", -1, "")
	log.Printf("Delete OK: %b", response.Ok)
	response, _ = core.Get(true, "twitter", "tweet", "1")
	log.Printf("Get: %s", response.Exists)

	healthResponse, _ := cluster.Health(true)
	log.Printf("Health: %s", healthResponse.Status)

	cluster.State("transient", "discovery.zen.minimum_master_nodes", 2)

}
