package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"log"
)

// The get API allows to get a typed JSON document from the index based on its id.
// GET - retrieves the doc
// HEAD - checks for existence of the doc
// http://www.elasticsearch.org/guide/reference/api/get.html
// TODO: make this implement an interface
func RunGet(pretty bool, index string, _type string, id string) (api.BaseResponse, error) {
	var response map[string]interface{}
	var body string
	var url string
	var retval api.BaseResponse

	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, api.Pretty(pretty))
	} else {
		url = fmt.Sprintf("/%s/%s?%s", index, id, api.Pretty(pretty))
	}
	req, err := api.ElasticSearchRequest("GET", url)
	if err != nil {
		// some sort of generic error handler		
	}
	body, err = req.Do(&response)
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	} else {
		// marshall into json
		jsonErr := json.Unmarshal([]byte(body), &retval)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	}
	fmt.Println(body)
	return retval, err
}
