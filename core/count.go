package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

type CountResponse struct {
	Count int        `json:"count"`
	Shard api.Status `json:"_shards"`
}

// The count API allows to easily execute a query and get the number of matches for that query. 
// It can be executed across one or more indices and across one or more types. 
// The query can either be provided using a simple query string as a parameter, 
//or using the Query DSL defined within the request body.
// http://www.elasticsearch.org/guide/reference/api/count.html
// TODO: take parameters. 
// currently not working against 0.19.10
func Count(pretty bool, index string, _type string) (CountResponse, error) {
	var url string
	var retval CountResponse
	url = fmt.Sprintf("/%s/%s/_count?%s", index, _type, api.Pretty(pretty))
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
