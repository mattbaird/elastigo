package indices

import (
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"log"
)

// Lists status details of all indices or the specified index.
// http://www.elasticsearch.org/guide/reference/api/admin-indices-status.html
func RunStatus(pretty bool, indices ...string) {
	index := ""
	if len(indices) >= 1 {
		index = indices[0]
	}

	var response map[string]interface{}

	var body string
	if len(index) > 0 {
		body = api.ElasticSearchRequest("GET", "/"+index+"/_status?pretty=1").Do(&response)
	} else {
		body = api.ElasticSearchRequest("GET", "/_status?pretty=1").Do(&response)
	}

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
