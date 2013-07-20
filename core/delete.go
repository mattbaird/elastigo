package core

import (
	"encoding/json"
	"fmt"
	"github.com/meanpath/elastigo/api"
)

// The delete API allows to delete a typed JSON document from a specific index based on its id. 
// http://www.elasticsearch.org/guide/reference/api/delete.html
// todo: add routing and versioning support
func Delete(pretty bool, index string, _type string, id string, version int, routing string) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, api.Pretty(pretty))
	body, err := api.DoCommand("DELETE", url, nil)
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
