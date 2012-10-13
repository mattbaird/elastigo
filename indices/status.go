package indices

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// Lists status details of all indices or the specified index.
// http://www.elasticsearch.org/guide/reference/api/admin-indices-status.html
func Status(pretty bool, indices ...string) (api.BaseResponse, error) {
	var retval api.BaseResponse
	var body string
	var url string
	if len(indices) > 0 {
		//TODO, emit the indices csv style

		url = "/" + indices[0] + "/_status?pretty=1"
	} else {
		url = "/_status?pretty=1"
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
