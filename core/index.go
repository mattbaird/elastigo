package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// The index API adds or updates a typed JSON document in a specific index, making it searchable.
// http://www.elasticsearch.org/guide/reference/api/index_.html
func Index(pretty bool, index string, _type string, id string, data interface{}) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, api.Pretty(pretty))
	var method string
	if id == "" {
		method = "POST"
	} else {
		method = "PUT"
	}

	body, err := api.DoCommand(method, url, data)
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
	//fmt.Println(body)
	return retval, err
}
