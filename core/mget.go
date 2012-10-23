package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// Multi GET API allows to get multiple documents based on an index, type (optional) and id (and possibly routing).
// The response includes a docs array with all the fetched documents, each element similar in structure to a document 
// provided by the get API.
// see http://www.elasticsearch.org/guide/reference/api/multi-get.html
func MGet(pretty bool, index string, _type string, mgetRequest MGetRequestContainer) (MGetResponseContainer, error) {
	var url string
	var retval MGetResponseContainer
	if len(index) <= 0 {
		url = fmt.Sprintf("/_mget?%s", api.Pretty(pretty))
	}
	if len(_type) > 0 && len(index) > 0 {
		url = fmt.Sprintf("/%s/%s/_mget?%s", index, _type, api.Pretty(pretty))
	} else if len(index) > 0 {
		url = fmt.Sprintf("/%s/_mget?%s", index, api.Pretty(pretty))
	}
	body, err := api.DoCommand("GET", url, nil)
	if err != nil {
		return retval, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal([]byte(body), &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}
	fmt.Println(body)
	return retval, err
}

type MGetRequestContainer struct {
	Docs []MGetRequest `json:"docs"`
}

type MGetRequest struct {
	Index  string   `json:"_index"`
	Type   string   `json:"_type"`
	ID     string   `json:"_id"`
	IDS    []string `json:"_ids,omitifempty"`
	Fields []string `json:"fields,omitifempty"`
}

type MGetResponseContainer struct {
	Docs []api.BaseResponse `json:"docs"`
}
