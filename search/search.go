package search

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"log"
	"net/url"
	"strings"
)

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

func (s *SearchDsl) Result() (*core.SearchResult, error) {
	var retval core.SearchResult
	if core.DebugRequests {
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
