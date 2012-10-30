package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// Performs a very basic search on an index via the request URI API.
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func Search(pretty bool, index string, _type string, query interface{}, scroll string) (SearchResult, error) {
	var url string
	var retval SearchResult
	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/_search?%s%s", index, _type, api.Pretty(pretty), api.Scroll(scroll))
	} else {
		url = fmt.Sprintf("/%s/_search?%s%s", index, api.Pretty(pretty), api.Scroll(scroll))
	}
	body, err := api.DoCommand("POST", url, query)
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
	return retval, err
}

func Scroll(pretty bool, scroll_id string, scroll string) (SearchResult, error) {
	var url string
	var retval SearchResult

	url = fmt.Sprintf("/_search/scroll?%s%s", api.Pretty(pretty), api.Scroll(scroll))

	body, err := api.DoCommand("POST", url, scroll_id)
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
	return retval, err
}

type SearchResult struct {
	Took        int        `json:"took"`
	TimedOut    bool       `json:"timed_out"`
	ShardStatus api.Status `json:"_shards"`
	Hits        Hits       `json:"hits"`
	Facets      json.RawMessage     `json:"facets,omitempty"` // structure varies on query
	ScrollId    string     `json:"_scroll_id,omitempty"`
}

type Hits struct {
	Total int `json:"total"`
	//	MaxScore float32 `json:"max_score"`
	Hits []Hit `json:"hits"`
}
type Hit struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type,omitempty"`
	Id     string          `json:"_id"`
	Score  float32         `json:"score"`
	Source json.RawMessage `json:"_source"` // marshalling left to consumer
}
