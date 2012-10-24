package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"net/url"
)

// Performs a very basic search on an index via the request URI API.
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func Search(pretty bool, index string, searchValues ...string) (api.BaseResponse, error) {
	query := searchValues[0]
	values := url.Values{}
	values.Set("q", query)
	var response map[string]interface{}
	var retval api.BaseResponse
	url := fmt.Sprintf("/%s/_search?%s%s", index, api.Pretty(pretty), values.Encode())
	req, err := api.ElasticSearchRequest("GET", url)
	if err != nil {
		return retval, err
	}
	body, err := req.Do(&response)
	if err != nil {
		return retval, err
	}
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		return retval, errors.New(fmt.Sprintf("Error: %v (%v)\n", error, status))
	} else {
		// marshall into json
		jsonErr := json.Unmarshal([]byte(body), &retval)
		if jsonErr != nil {
			return retval, err
		}
	}
	return retval, err
}
