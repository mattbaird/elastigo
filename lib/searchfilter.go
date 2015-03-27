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

	. "github.com/araddon/gou"
)

var (
	_ = DEBUG
)

// A bool (and/or) clause
type BoolClause string

// Filter clause is either a boolClause or FilterOp
type FilterClause interface {
	String() string
}

// A wrapper to allow for custom serialization
type FilterWrap struct {
	boolClause string
	filters    []interface{}
}

func NewFilterWrap() *FilterWrap {
	return &FilterWrap{filters: make([]interface{}, 0), boolClause: "and"}
}

func (f *FilterWrap) String() string {
	return fmt.Sprintf(`fopv: %d:%v`, len(f.filters), f.filters)
}

// Custom marshalling to support the query dsl
func (f *FilterWrap) addFilters(fl []interface{}) {
	if len(fl) > 1 {
		fc := fl[0]
		switch fc.(type) {
		case BoolClause, string:
			f.boolClause = fc.(string)
			fl = fl[1:]
		}
	}
	f.filters = append(f.filters, fl...)
}

// Custom marshalling to support the query dsl
func (f *FilterWrap) MarshalJSON() ([]byte, error) {
	var root interface{}
	if len(f.filters) > 1 {
		root = map[string]interface{}{f.boolClause: f.filters}
	} else if len(f.filters) == 1 {
		root = f.filters[0]
	}
	return json.Marshal(root)
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

	"filter" : {
	    "terms" : {
	        "user" : ["kimchy", "elasticsearch"],
	        "execution" : "bool",
	        "_cache": true
	    }
	}

	"filter" : {
	    "term" : { "user" : "kimchy"}
	}

	"filter" : {
	    "and" : [
	        {
	            "range" : {
	                "postDate" : {
	                    "from" : "2010-03-01",
	                    "to" : "2010-04-01"
	                }
	            }
	        },
	        {
	            "prefix" : { "name.second" : "ba" }
	        }
	    ]
	}

*/

// Filter Operation
//
//   Filter().Term("user","kimchy")
//
//   // we use variadics to allow n arguments, first is the "field" rest are values
//   Filter().Terms("user", "kimchy", "elasticsearch")
//
//   Filter().Exists("repository.name")
//
func Filter() *FilterOp {
	return &FilterOp{}
}

func CompoundFilter(fl ...interface{}) *FilterWrap {
	FilterVal := NewFilterWrap()
	FilterVal.addFilters(fl)
	return FilterVal
}

type FilterOp struct {
	curField    string
	TermsMap    map[string][]interface{} `json:"terms,omitempty"`
	TermMap     map[string]interface{}   `json:"term,omitempty"`
	RangeMap    map[string]RangeFilter   `json:"range,omitempty"`
	Exist       map[string]string        `json:"exists,omitempty"`
	MisssingVal map[string]string        `json:"missing,omitempty"`
	AndFilters  []FilterOp               `json:"and,omitempty"`
	OrFilters   []FilterOp               `json:"or,omitempty"`
	Limit       *LimitFilter             `json:"limit,omitempty"`
}

type LimitFilter struct {
	Value int `json:"value"`
}

type RangeFilter struct {
	Gte      interface{} `json:"gte,omitempty"`
	Lte      interface{} `json:"lte,omitempty"`
	Gt       interface{} `json:"gt,omitempty"`
	Lt       interface{} `json:"lt,omitempty"`
	TimeZone string      `json:"time_zone,omitempty"` //Ideally this would be an int
}

// A range is a special type of Filter operation
//
//    Range().Exists("repository.name")
func Range() *FilterOp {
	return &FilterOp{RangeMap: make(map[string]RangeFilter)}
}

// Term will add a term to the filter.
// Multiple Term filters can be added, and ES will OR them.
func (f *FilterOp) Term(field string, value interface{}) *FilterOp {
	if len(f.TermMap) == 0 {
		f.TermMap = make(map[string]interface{})
	}

	f.TermMap[field] = value
	return f
}

func (f *FilterOp) And(filter *FilterOp) *FilterOp {
	if len(f.AndFilters) == 0 {
		f.AndFilters = []FilterOp{*filter}
	} else {
		f.AndFilters = append(f.AndFilters, *filter)
	}

	return f
}

func (f *FilterOp) Or(filter *FilterOp) *FilterOp {
	if len(f.OrFilters) == 0 {
		f.OrFilters = []FilterOp{*filter}
	} else {
		f.OrFilters = append(f.OrFilters, *filter)
	}

	return f
}

// Filter Terms
//
//   Filter().Terms("user","kimchy","stuff")
//	 Note: you can only have one terms clause in a filter. Use a bool filter to combine
func (f *FilterOp) Terms(field string, values ...interface{}) *FilterOp {
	//You can only have one terms in a filter
	f.TermsMap = make(map[string][]interface{})

	for _, val := range values {
		f.TermsMap[field] = append(f.TermsMap[field], val)
	}

	return f
}

// AddRange adds a range filter for the given field.
func (f *FilterOp) AddRange(field string, gte interface{},
	gt interface{}, lte interface{}, lt interface{}, timeZone string) *FilterOp {

	if f.RangeMap == nil {
		f.RangeMap = make(map[string]RangeFilter)
	}

	f.RangeMap[field] = RangeFilter{
		Gte:      gte,
		Gt:       gt,
		Lte:      lte,
		Lt:       lt,
		TimeZone: timeZone}

	return f
}

func (f *FilterOp) SetLimit(maxResults int) *FilterOp {
	f.Limit = &LimitFilter{Value: maxResults}
	return f
}
