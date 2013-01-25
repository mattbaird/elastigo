package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"log"
	"net/url"
	"strings"
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
	log.Println(uriVal)
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
	Score  float32         `json:"score"`
	Source json.RawMessage `json:"_source"` // marshalling left to consumer
}

// This is the entry point to the SearchDsl, it is a chainable set of utilities
// to create searches.   
//
// params
//    @index = elasticsearch index to search
// 
//    out, err := Search("github").Type("Issues").Pretty().Query(
//    Query().Range(
//         Range().Field("created_at").From("2012-12-10T15:00:00-08:00").To("2012-12-10T15:10:00-08:00"),
//       ).Search("add"),
//     ).Result()
func Search(index string) *SearchDsl {
	return &SearchDsl{Index: index, args: url.Values{}}
}

type SearchDsl struct {
	args      url.Values
	types     []string
	FromVal   int        `json:"from,omitempty"`
	SizeVal   int        `json:"size,omitempty"`
	Index     string     `json:"-"`
	FacetVal  *FacetDsl  `json:"facets,omitempty"`
	QueryVal  *QueryDsl  `json:"query,omitempty"`
	SortBody  []*SortDsl `json:"sort,omitempty"`
	FilterVal *FilterOp  `json:"filter,omitempty"`
}

func (s *SearchDsl) Bytes() ([]byte, error) {

	return api.DoCommand("POST", s.url(), s)
}

func (s *SearchDsl) Result() (*SearchResult, error) {
	var retval SearchResult
	if DebugRequests {
		sb, _ := json.MarshalIndent(s, "  ", "  ")
		log.Println(s.url())
		log.Println(string(sb))
	}
	body, err := s.Bytes()
	if err != nil {
		return nil, err
	}
	jsonErr := json.Unmarshal([]byte(body), &retval)
	return &retval, jsonErr
}

func (s *SearchDsl) url() string {
	url := fmt.Sprintf("/%s%s/_search?%s", s.Index, s.getType(), s.args.Encode())
	return url
}

func (s *SearchDsl) Pretty() *SearchDsl {
	s.args.Set("pretty", "1")
	return s
}

// this is the elasticsearch *Type* within a specific index
func (s *SearchDsl) Type(indexType string) *SearchDsl {
	if len(s.types) == 0 {
		s.types = make([]string, 0)
	}
	s.types = append(s.types, indexType)
	return s
}

func (s *SearchDsl) getType() string {
	if len(s.types) > 0 {
		return "/" + strings.Join(s.types, ",")
	}
	return ""
}

func (s *SearchDsl) From(from string) *SearchDsl {
	s.args.Set("from", from)
	return s
}

// This is a simple interfaceto search, doesn't have the power of query
// but uses a simple query_string search
func (s *SearchDsl) Search(srch string) *SearchDsl {
	s.QueryVal = Query().Search(srch)
	return s
}

func (s *SearchDsl) Size(size string) *SearchDsl {
	s.args.Set("size", size)
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
func (s *SearchDsl) Filter(f *FilterOp) *SearchDsl {
	s.FilterVal = f
	return s
}

func (s *SearchDsl) Sort(sort ...*SortDsl) *SearchDsl {
	if s.SortBody == nil {
		s.SortBody = make([]*SortDsl, 0)
	}
	s.SortBody = append(s.SortBody, sort...)
	return s
}

/* 
	Sorting accepts any number of Sort commands

	Query().Sort(
		Sort("last_name").Desc(),
		Sort("age"),
	)
*/
func Sort(field string) *SortDsl {
	return &SortDsl{Name: field}
}

type SortBody []interface{}
type SortDsl struct {
	Name   string
	IsDesc bool
}

func (s *SortDsl) Desc() *SortDsl {
	s.IsDesc = true
	return s
}
func (s *SortDsl) Asc() *SortDsl {
	s.IsDesc = false
	return s
}
func (s *SortDsl) MarshalJSON() ([]byte, error) {
	if s.IsDesc {
		return json.Marshal(map[string]string{s.Name: "desc"})
	}
	if s.Name == "_score" {
		return []byte(`"_score"`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, s.Name)), nil // "user"  assuming default = asc?
	// TODO
	//    { "price" : {"missing" : "_last"} },
	//    { "price" : {"ignore_unmapped" : true} },
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
func Facet() *FacetDsl {
	return &FacetDsl{&Term{Terms{nil, ""}}}
}

type FacetDsl struct {
	TermsVal *Term `json:"terms,omitempty"`
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

// Generic Term based (used in query, facet, filter)
type Term struct {
	Terms Terms `json:"terms,omitempty"`
}

type Terms struct {
	Fields []string `json:"field,omitempty"`
	Size   string   `json:"size,omitempty"`
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
	      "query": " actor:\"bob\"  AND type:\"EventType\""
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
	Terms    map[string]string `json:"term,omitempty"`
	Qs       *QueryString      `json:"query_string,omitempty"`
	//Exist    string            `json:"_exists_,omitempty"`
}

// Custom marshalling
func (s *QueryDsl) MarshalJSON() ([]byte, error) {
	if s.IsDesc {
		return json.Marshal(map[string]string{s.Name: "desc"})
	}
	if s.Name == "_score" {
		return []byte(`"_score"`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, s.Name)), nil // "user"  assuming default = asc?
	// TODO
	//    { "price" : {"missing" : "_last"} },
	//    { "price" : {"ignore_unmapped" : true} },
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

// Add a term search for a specific field 
//    Term("user","kimchy")
func (q *QueryDsl) Term(name, value string) *QueryDsl {
	if len(q.Terms) == 0 {
		q.Terms = make(map[string]string)
	}
	q.Terms[name] = value
	return q
}
func (q *QueryDsl) Search(searchFor string) *QueryDsl {
	//I don't think this is right, it is not a filter.query, it should be q query?
	qs := NewQueryString()
	q.Qs = &qs
	q.Qs.Query = searchFor
	return q
}

// Fields in query_string search
//     Fields("fieldname","search_for","","")
//     
//     Fields("fieldname,field2,field3","search_for","","")
//
//     Fields("fieldname,field2,field3","search_for","field_exists","")
func (q *QueryDsl) Fields(fields, search, exists, missing string) *QueryDsl {
	fieldList := strings.Split(fields, ",")
	qs := NewQueryString()
	q.Qs = &qs
	q.Qs.Query = search
	if len(fieldList) == 1 {
		q.Qs.DefaultField = fields
	} else {
		q.Qs.Fields = fieldList
	}
	q.Qs.Exists = exists
	q.Qs.Missing = missing
	return q
}

type MatchAll struct {
	All string `json:"-"`
}

type Filtered struct {
	Query  *QueryWrap `json:"query,omitempty"`
	Filter *FilterOp  `json:"filter,omitempty"`
}

// should we reuse QueryDsl here?
type QueryWrap struct {
	Qs QueryString `json:"query_string,omitempty"`
}

/*
type ExistsString string

func (e *ExistsString) MarshallJSON() ([]byte, error) {
	if len(e) == 0 {
		return nil, nil
	}
	return []byte(e), nil
}
*/
func NewQueryString() QueryString {
	return QueryString{"", "", "", "", "", nil}
}

type QueryString struct {
	DefaultOperator string   `json:"default_operator,omitempty"`
	DefaultField    string   `json:"default_field,omitempty"`
	Query           string   `json:"query,omitempty"`
	Exists          string   `json:"_exists_,omitempty"`
	Missing         string   `json:"_missing_,omitempty"`
	Fields          []string `json:"fields,omitempty"`
	//_exists_:field1,
	//_missing_:field1,
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
"filter": {
    "missing": {
        "field": "repository.name"
    }
}

*/

func Filter() *FilterOp {
	return &FilterOp{}
}

type FilterOp struct {
	curField    string
	Range       map[string]map[string]string `json:"range,omitempty"`
	Exist       map[string]string            `json:"exists,omitempty"`
	MisssingVal map[string]string            `json:"missing,omitempty"`
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
func (f *FilterOp) Exists(name string) *FilterOp {
	f.Exist = map[string]string{"field": name}
	return f
}
func (f *FilterOp) Missing(name string) *FilterOp {
	f.MisssingVal = map[string]string{"field": name}
	return f
}
