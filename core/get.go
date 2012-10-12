package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// The get API allows to get a typed JSON document from the index based on its id.
// GET - retrieves the doc
// HEAD - checks for existence of the doc
// http://www.elasticsearch.org/guide/reference/api/get.html
// TODO: make this implement an interface
func Get(pretty bool, index string, _type string, id string) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, api.Pretty(pretty))
	} else {
		url = fmt.Sprintf("/%s/%s?%s", index, id, api.Pretty(pretty))
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

// The API also allows to check for the existance of a document using HEAD
// This appears to be broken in the current version of elasticsearch 0.19.10, currently
// returning nothing
func Exists(pretty bool, index string, _type string, id string) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, api.Pretty(pretty))
	} else {
		url = fmt.Sprintf("/%s/%s?%s", index, id, api.Pretty(pretty))
	}
	body, err := api.DoCommand("HEAD", url, nil)
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
