package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// The validate API allows a user to validate a potentially expensive query without executing it.
// see http://www.elasticsearch.org/guide/reference/api/validate.html
func Validate(pretty bool, index string, _type string, query string, explain bool) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/_validate/query?q=%s&%s&explain=%s", index, _type, query, api.Pretty(pretty), explain)
	} else {
		url = fmt.Sprintf("/%s/_validate/query?q=%s&%s&explain=%s", index, query, api.Pretty(pretty), explain)
	}
	body, err := api.DoCommand("GET", url, nil)
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

type Validation struct {
	Valid         bool           `json:"valid"`
	Shards        api.Status     `json:"_shards"`
	Explainations []Explaination `json:"explanations,omitempty"`
}

type Explaination struct {
	Index string `json:"index"`
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}
