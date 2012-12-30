package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// The percolator allows to register queries against an index, and then send percolate requests which include a doc, and 
// getting back the queries that match on that doc out of the set of registered queries.
// Think of it as the reverse operation of indexing and then searching. Instead of sending docs, indexing them, 
// and then running queries. One sends queries, registers them, and then sends docs and finds out which queries
// match that doc.
// see http://www.elasticsearch.org/guide/reference/api/percolate.html
func RegisterPercolate(pretty bool, index string, name string, query api.Query) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	url = fmt.Sprintf("/_percolator/%s/%s?%s", index, name, api.Pretty(pretty))
	body, err := api.DoCommand("PUT", url, query)
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

func Percolate(pretty bool, index string, _type string, name string, doc string) (api.Match, error) {
	var url string
	var retval api.Match
	url = fmt.Sprintf("/%s/%s/_percolate?%s", index, _type, api.Pretty(pretty))
	body, err := api.DoCommand("GET", url, doc)
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
