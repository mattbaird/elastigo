package indices

import (
	"encoding/json"
	"fmt"
	"github.com/meanpath/elastigo/api"
)

// The flush API allows to flush one or more indices through an API. The flush process of an index basically
// frees memory from the index by flushing data to the index storage and clearing the internal transaction 
// log. By default, ElasticSearch uses memory heuristics in order to automatically trigger flush operations 
// as required in order to clear memory.
// http://www.elasticsearch.org/guide/reference/api/admin-indices-flush.html
// TODO: add Shards to response
func Flush(index ...string) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	if len(index) > 0 {
		url = fmt.Sprintf("/%s/_flush", index)
	} else {
		url = "/_flush"
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
	//fmt.Println(body)
	return retval, err
}
