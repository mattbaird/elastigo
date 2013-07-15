package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"net/url"
	"strconv"
)

var (
	DebugRequests = false
)

// Performs a very basic search on an index via the request URI API.
//
// params:
//   @pretty:  bool for pretty reply or not, a parameter to elasticsearch
//   @index:  the elasticsearch index
//   @_type:  optional ("" if not used) search specific type in this index
//   @query:  this can be one of 3 types:
//              1)  string value that is valid elasticsearch
//              2)  io.Reader that can be set in body (also valid elasticsearch string syntax..)
//              3)  other type marshalable to json (also valid elasticsearch json)
//
//   out, err := SearchRequest(true, "github","",qryType ,"")
//
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func SearchRequest(pretty bool, index string, _type string, query interface{}, scroll string) (SearchResult, error) {
	var uriVal string
	var retval SearchResult
	if len(_type) > 0 && _type != "*" {
		uriVal = fmt.Sprintf("/%s/%s/_search?%s%s", index, _type, api.Pretty(pretty), api.Scroll(scroll))
	} else {
		uriVal = fmt.Sprintf("/%s/_search?%s%s", index, api.Pretty(pretty), api.Scroll(scroll))
	}
	body, err := api.DoCommand("POST", uriVal, query)
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

// Performs the simplest possible query in url string
// params:
//   @index:  the elasticsearch index
//   @_type:  optional ("" if not used) search specific type in this index
//   @query:  valid string lucene search syntax
//
//   out, err := SearchUri("github","",`user:kimchy` ,"")
//
// produces a request like this:    host:9200/github/_search?q=user:kimchy"
//
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func SearchUri(index, _type string, query, scroll string) (SearchResult, error) {
	var uriVal string
	var retval SearchResult
	query = url.QueryEscape(query)
	if len(_type) > 0 && _type != "*" {
		uriVal = fmt.Sprintf("/%s/%s/_search?q=%s%s", index, _type, query, api.Scroll(scroll))
	} else {
		uriVal = fmt.Sprintf("/%s/_search?q=%s%s", index, query, api.Scroll(scroll))
	}
	//log.Println(uriVal)
	body, err := api.DoCommand("GET", uriVal, nil)
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
	Took        int             `json:"took"`
	TimedOut    bool            `json:"timed_out"`
	ShardStatus api.Status      `json:"_shards"`
	Hits        Hits            `json:"hits"`
	Facets      json.RawMessage `json:"facets,omitempty"` // structure varies on query
	ScrollId    string          `json:"_scroll_id,omitempty"`
}

func (s *SearchResult) String() string {
	return fmt.Sprintf("<Results took=%v Timeout=%v hitct=%v />", s.Took, s.TimedOut, s.Hits.Total)
}

type Hits struct {
	Total int `json:"total"`
	//	MaxScore float32 `json:"max_score"`
	Hits []Hit `json:"hits"`
}

func (h *Hits) Len() int {
	return len(h.Hits)
}

type Hit struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type,omitempty"`
	Id     string          `json:"_id"`
	Score  Float32Nullable `json:"_score,omitempty"` // Filters (no query) dont have score, so is null
	Source json.RawMessage `json:"_source"`          // marshalling left to consumer
}

// Elasticsearch returns some invalid (according to go) json, with floats having...
//
// json: cannot unmarshal null into Go value of type float32 (see last field.)
//
// "hits":{"total":6808,"max_score":null,
//    "hits":[{"_index":"10user","_type":"user","_id":"751820","_score":null,
type Float32Nullable float32

func (i *Float32Nullable) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	if in, err := strconv.ParseFloat(string(data), 32); err != nil {
		return err
	} else {
		*i = Float32Nullable(in)
	}
	return nil
}
