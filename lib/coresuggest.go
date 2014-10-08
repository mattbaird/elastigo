// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elastigo

import (
	"encoding/json"
	"fmt"
	//"strconv"
	//"strings"
)

// Search performs a very basic search on an index via the request URI API.
//
// params:
//   @index:  the elasticsearch index
//   @_type:  optional ("" if not used) search specific type in this index
//   @args:   a map of URL parameters. Allows all the URI-request parameters allowed by ElasticSearch.
//   @query:  this can be one of 3 types:
//              1)  string value that is valid elasticsearch
//              2)  io.Reader that can be set in body (also valid elasticsearch string syntax..)
//              3)  other type marshalable to json (also valid elasticsearch json)
//
//   out, err := Search(true, "github", map[string]interface{} {"from" : 10}, qryType)
//
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func (c *Conn) Suggest(index string, _type string, args map[string]interface{}, query interface{}) (SuggestResult, error) {
	var uriVal string
	var retval SuggestResult
	if len(_type) > 0 && _type != "*" {
		uriVal = fmt.Sprintf("/%s/%s/_suggest", index, _type)
	} else {
		uriVal = fmt.Sprintf("/%s/_suggest", index)
	}
	body, err := c.DoCommand("POST", uriVal, args, query)
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

type SuggestResult struct {
	Suggestions []Suggestion `json:"suggestions"`
	ShardStatus Status       `json:"_shards"`
}

/*

type SuggestionOption struct {
	Payload json.RawMessage `json:"payload"`
	Score   Float32Nullable `json:"score,omitempty"`
	Text    string          `json:"text"`
}

type Suggestion struct {
	Length  int                `json:"length"`
	Offset  int                `json:"offset"`
	Options []SuggestionOption `json:"options"`
	Text    string             `json:"text"`
}

type Suggestions map[string][]Suggestion

type SearchResult struct {
	RawJSON      []byte
	Took         int             `json:"took"`
	TimedOut     bool            `json:"timed_out"`
	ShardStatus  Status          `json:"_shards"`
	Hits         Hits            `json:"hits"`
	Facets       json.RawMessage `json:"facets,omitempty"` // structure varies on query
	ScrollId     string          `json:"_scroll_id,omitempty"`
	Aggregations json.RawMessage `json:"aggregations,omitempty"` // structure varies on query
	Suggestions  Suggestions     `json:"suggest,omitempty"`
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

type Highlight map[string][]string

type Hit struct {
	Index       string           `json:"_index"`
	Type        string           `json:"_type,omitempty"`
	Id          string           `json:"_id"`
	Score       Float32Nullable  `json:"_score,omitempty"` // Filters (no query) dont have score, so is null
	Source      *json.RawMessage `json:"_source"`          // marshalling left to consumer
	Fields      *json.RawMessage `json:"fields"`           // when a field arg is passed to ES, instead of _source it returns fields
	Explanation *Explanation     `json:"_explanation,omitempty"`
	Highlight   *Highlight       `json:"highlight,omitempty"`
}

func (e *Explanation) String(indent string) string {
	if len(e.Details) == 0 {
		return fmt.Sprintf("%s>>>  %v = %s", indent, e.Value, strings.Replace(e.Description, "\n", "", -1))
	} else {
		detailStrs := make([]string, 0)
		for _, detail := range e.Details {
			detailStrs = append(detailStrs, fmt.Sprintf("%s", detail.String(indent+"| ")))
		}
		return fmt.Sprintf("%s%v = %s(\n%s\n%s)", indent, e.Value, strings.Replace(e.Description, "\n", "", -1), strings.Join(detailStrs, "\n"), indent)
	}
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

*/
