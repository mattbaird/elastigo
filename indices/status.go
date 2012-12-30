package indices

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"strings"
)

// Lists status details of all indices or the specified index.
// http://www.elasticsearch.org/guide/reference/api/admin-indices-status.html
func Status(pretty bool, indices ...string) (api.BaseResponse, error) {
	var retval api.BaseResponse
	var url string
	if len(indices) > 0 {
		url = fmt.Sprintf("/%s/_status?%s", strings.Join(indices, ","), api.Pretty(pretty))

	} else {
		url = fmt.Sprintf("/_status?%s", api.Pretty(pretty))
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
