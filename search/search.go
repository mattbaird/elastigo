package search

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"log"
	"net/url"
	"strings"

	. "github.com/araddon/gou"
)

var (
	_ = DEBUG
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
	FromVal   int         `json:"from,omitempty"`
	SizeVal   int         `json:"size,omitempty"`
	Index     string      `json:"-"`
	FacetVal  *FacetDsl   `json:"facets,omitempty"`
	QueryVal  *QueryDsl   `json:"query,omitempty"`
	SortBody  []*SortDsl  `json:"sort,omitempty"`
	FilterVal *FilterWrap `json:"filter,omitempty"`
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
		Logf(ERROR, "%v", err)
		return nil, err
	}
	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		Logf(ERROR, "%v \n\t%s", jsonErr, string(body))
	}
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

// Pass a Query expression to this search
//
//		qry := Search("github").Size("0").Facet(
//					Facet().Regex("repository.name", "no.*").Size("8"),
//				)
//
//		qry := Search("github").Pretty().Facet(
//					Facet().Fields("type").Size("25"),
//				)
func (s *SearchDsl) Facet(f *FacetDsl) *SearchDsl {
	s.FacetVal = f
	return s
}

func (s *SearchDsl) Query(q *QueryDsl) *SearchDsl {
	s.QueryVal = q
	return s
}

// Add Filter Clause with optional Boolean Clause.  This accepts n number of
// filter clauses.  If more than one, and missing Boolean Clause it assumes "and"
//
//     qry := Search("github").Filter(
//         Filter().Exists("repository.name"),
//     )
//
//     qry := Search("github").Filter(
//         "or",
//         Filter().Exists("repository.name"),
//         Filter().Terms("actor_attributes.location", "portland"),
//     )
//
//     qry := Search("github").Filter(
//         Filter().Exists("repository.name"),
//         Filter().Terms("repository.has_wiki", true)
//     )
func (s *SearchDsl) Filter(fl ...interface{}) *SearchDsl {
	if s.FilterVal == nil {
		s.FilterVal = NewFilterWrap()
	}

	s.FilterVal.addFilters(fl)
	return s
}

func (s *SearchDsl) Sort(sort ...*SortDsl) *SearchDsl {
	if s.SortBody == nil {
		s.SortBody = make([]*SortDsl, 0)
	}
	s.SortBody = append(s.SortBody, sort...)
	return s
}
