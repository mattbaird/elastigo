package core

import (
	"encoding/json"
	"fmt"
	"github.com/mschoch/elastigo/api"
	"strings"
)

// The delete by query API allows to delete documents from one or more indices and one or more types based on a query. 
// The query can either be provided using a simple query string as a parameter, or using the Query DSL defined within 
// the request body.
// see: http://www.elasticsearch.org/guide/reference/api/delete-by-query.html
func DeleteByQuery(pretty bool, indices []string, types []string, query interface{}) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	if len(indices) > 0 && len(types) > 0 {
		url = fmt.Sprintf("http://localhost:9200/%s/%s/_query?%s&%s", strings.Join(indices, ","), strings.Join(types, ","), buildQuery, api.Pretty(pretty))
	} else if len(indices) > 0 {
		url = fmt.Sprintf("http://localhost:9200/%s/_query?%s&%s", strings.Join(indices, ","), buildQuery, api.Pretty(pretty))
	}
	body, err := api.DoCommand("DELETE", url, query)
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

func buildQuery() string {
	return ""
}

type DeleteByQueryResponse struct {
	Status   bool                   `json:"ok"`
	Indicies map[string]IndexStatus `json:"_indices"`
}

type IndexStatus struct {
	Shards api.Status `json:"_shards"`
}
