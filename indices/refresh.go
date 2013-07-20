package indices

import (
	"encoding/json"
	"fmt"
	"github.com/meanpath/elastigo/api"
	"strings"
)

// The refresh API allows to explicitly refresh one or more index, making all operations performed since 
// the last refresh available for search. The (near) real-time capabilities depend on the index engine used. 
// For example, the robin one requires refresh to be called, but by default a refresh is scheduled periodically.
// http://www.elasticsearch.org/guide/reference/api/admin-indices-refresh.html
// TODO: add Shards to response
func Refresh(indices ...string) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	if len(indices) > 0 {
		url = fmt.Sprintf("/%s/_refresh", strings.Join(indices, ","))
	} else {
		url = "/_refresh"
	}
	body, err := api.DoCommand("POST", url, nil)
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
	return retval, err
}
