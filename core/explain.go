package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// The explain api computes a score explanation for a query and a specific document. 
// This can give useful feedback whether a document matches or didn’t match a specific query. 
// This feature is available from version 0.19.9 and up.
// see http://www.elasticsearch.org/guide/reference/api/explain.html
func Explain(pretty bool, index string, _type string, id string, query string) (api.Match, error) {
	var url string
	var retval api.Match
	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/_explain?%s", index, _type, api.Pretty(pretty))
	} else {
		url = fmt.Sprintf("/%s/_explain?%s", index, api.Pretty(pretty))
	}
	body, err := api.DoCommand("GET", url, query)
	if err != nil {
		return retval, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}
	fmt.Println(body)
	return retval, err
}
