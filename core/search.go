package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"log"
	"net/url"
)

// Performs a very basic search on an index via the request URI API.
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func RunSearch(pretty bool, index string, searchValues ...string) {
	query := searchValues[0]
	values := url.Values{}
	values.Set("q", query)
	var response map[string]interface{}

	url := fmt.Sprintf("/%s/_search?%s%s", index, api.Pretty(pretty), values.Encode())
	req, err := api.ElasticSearchRequest("GET", url)
	body, err := req.Do(&response)
	if err != nil {
		// some sort of generic error handler		
	}

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	} else {
		// marshall into json
		var objResponse api.BaseResponse
		jsonErr := json.Unmarshal([]byte(body), &objResponse)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	}
	fmt.Print(body)
}
