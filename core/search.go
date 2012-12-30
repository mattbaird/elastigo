package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"log"
	"strings"
)

// Performs a very basic search on an index via the request URI API.
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func SearchRequest(pretty bool, index string, _type string, query interface{}, scroll, size string) (SearchResult, error) {
	var url string
	var retval SearchResult
	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/_search?%s%s%s", index, _type, api.Pretty(pretty), api.Scroll(scroll), api.Size(size))
	} else {
		url = fmt.Sprintf("/%s/_search?%s%s%s", index, api.Pretty(pretty), api.Scroll(scroll), api.Size(size))
	}
	log.Println(url)
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
	Took        int             `json:"took"`
	TimedOut    bool            `json:"timed_out"`
	ShardStatus api.Status      `json:"_shards"`
	Hits        Hits            `json:"hits"`
	Facets      json.RawMessage `json:"facets,omitempty"` // structure varies on query
	ScrollId    string          `json:"_scroll_id,omitempty"`
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

func Search(index string) *SearchDsl {
	return &SearchDsl{Index: index, args: make([]string, 0)}
}

type SearchDsl struct {
	args      []string
	FromVal   int       `json:"from,omitempty"`
	SizeVal   int       `json:"size,omitempty"`
	Index     string    `json:"-"`
	IndexType string    `json:"-"`
	FacetVal  *FacetDsl `json:"facets"`
	QueryVal  *QueryDsl `json:"query,omitempty"`
	//FilterVal    FilterDsl `json:"filter,omitempty"`
}

func (s *SearchDsl) Bytes() ([]byte, error) {

	return api.DoCommand("POST", s.url(), s)
}

func (s *SearchDsl) Result() (*SearchResult, error) {
	var retval SearchResult
	sb, _ := json.MarshalIndent(s, "  ", "  ")
	log.Println(string(sb))
	body, err := s.Bytes()
	if err != nil {
		return nil, err
	}
	jsonErr := json.Unmarshal([]byte(body), &retval)
	return &retval, jsonErr
}

func (s *SearchDsl) url() string {
	url := fmt.Sprintf("/%s%s/_search?%s", s.Index, s.IndexType, strings.Join(s.args, "&"))
	log.Println(url)
	return url
}
func (s *SearchDsl) Pretty() *SearchDsl {
	s.args = append(s.args, "pretty=1")
	return s
}
func (s *SearchDsl) Type(indexType string) *SearchDsl {
	s.IndexType = "/" + indexType
	return s
}
func (s *SearchDsl) From(from string) *SearchDsl {
	s.args = append(s.args, "from="+from)
	return s
}
func (s *SearchDsl) Size(size string) *SearchDsl {
	s.args = append(s.args, "size="+size)
	return s
}
func (s *SearchDsl) Facet(f *FacetDsl) *SearchDsl {
	s.FacetVal = f
	return s
}
func (s *SearchDsl) Query(q *QueryDsl) *SearchDsl {
	s.QueryVal = q
	return s
}

func Facet() *FacetDsl {
	return &FacetDsl{&FacetTerm{FacetTerms{nil, ""}}}
}

/*
"facets": {
    "terms": {
		"terms": {
			"field": [
			  "@fields.category"
			],
			"size": 25
		}
    }
}
*/
type FacetDsl struct {
	TermsVal *FacetTerm `json:"terms,omitempty"`
}
type FacetTerm struct {
	Terms FacetTerms `json:"terms,omitempty"`
}

type FacetTerms struct {
	Fields []string `json:"field,omitempty"`
	Size   string   `json:"size,omitempty"`
}

func (f *FacetDsl) Size(size string) *FacetDsl {
	f.TermsVal.Terms.Size = size
	return f
}
func (f *FacetDsl) Fields(fields ...string) *FacetDsl {
	flds := make([]string, 0)
	for _, field := range fields {
		flds = append(flds, field)
	}
	f.TermsVal.Terms.Fields = flds
	return f
}

func Query() *QueryDsl {
	return &QueryDsl{}
}

/*

Three ways to serialize this query term
"query": {
	"filtered": {
	  "query": {
	    "query_string": {
	      "default_operator": "OR",
	      "default_field": "_all",
	      "query": " @fields.aid:\"10\"  AND @fields.PageType:\"*\""
	    }
	  },
	  "filter": {
	    "range": {
	      "@timestamp": {
	        "from": "2012-12-29T16:52:48+00:00",
	        "to": "2012-12-29T17:52:48+00:00"
	      }
	    }
	  }
	}
},

"query" : {
    "term" : { "user" : "kimchy" }
}

"query" : {
    "match_all" : {}
},
*/
type QueryDsl struct {
	Filter   *Filtered         `json:"filtered,omitempty"`
	MatchAll *MatchAll         `json:"match_all,omitempty"`
	Term     map[string]string `json:"term,omitempty"`
}

func (q *QueryDsl) All() *QueryDsl {
	q.MatchAll = &MatchAll{""}
	return q
}
func (q *QueryDsl) Range(fop *FilterOp) *QueryDsl {
	if q.Filter == nil {
		q.Filter = &Filtered{nil, nil}
	}
	q.Filter.Filter = fop
	return q
}
func (q *QueryDsl) Search(qs string) *QueryDsl {
	if q.Filter == nil {
		q.Filter = &Filtered{nil, nil}
	}
	q.Filter.Query = &FilterQueryWrap{FilterQuery{"", "", qs}}
	return q
}

type MatchAll struct {
	All string `json:"-"`
}

type Filtered struct {
	Query  *FilterQueryWrap `json:"query,omitempty"`
	Filter *FilterOp        `json:"filter,omitempty"`
}
type FilterQueryWrap struct {
	Query FilterQuery `json:"query_string,omitempty"`
}
type FilterQuery struct {
	DefaultOperator string `json:"default_operator,omitempty"`
	DefaultField    string `json:"default_field,omitempty"`
	Query           string `json:"query,omitempty"`
}

/*
"filter": {
	"range": {
	  "@timestamp": {
	    "from": "2012-12-29T16:52:48+00:00",
	    "to": "2012-12-29T17:52:48+00:00"
	  }
	}
}
*/
type FilterOp struct {
	curField string
	Range    map[string]map[string]string `json:"range,omitempty"`
}

func Range() *FilterOp {
	return &FilterOp{Range: make(map[string]map[string]string)}
}

func (f *FilterOp) Field(fld string) *FilterOp {
	f.curField = fld
	if _, ok := f.Range[fld]; !ok {
		m := make(map[string]string)
		f.Range[fld] = m
	}
	return f
}
func (f *FilterOp) From(from string) *FilterOp {
	f.Range[f.curField]["from"] = from
	return f
}
func (f *FilterOp) To(to string) *FilterOp {
	f.Range[f.curField]["to"] = to
	return f
}
