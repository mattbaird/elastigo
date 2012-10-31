package core

import (
	"encoding/json"
	"fmt"
	"github.com/mschoch/elastigo/api"
)

// The index API adds or updates a typed JSON document in a specific index, making it searchable. 
// http://www.elasticsearch.org/guide/reference/api/index_.html
func Index(pretty bool, index string, _type string, id string, data interface{}) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, api.Pretty(pretty))
	body, err := api.DoCommand("PUT", url, data)
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
