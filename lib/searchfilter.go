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

type TermExecutionMode string

const (
	TEM_DEFAULT TermExecutionMode = ""
	TEM_PLAIN                     = "plain"
	TEM_FIELD                     = "field_data"
	TEM_BOOL                      = "bool"
	TEM_AND                       = "and"
	TEM_OR                        = "or"
)

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

// Bool sets the type of boolean filter to use.
// Accepted values are "and" and "or".
func (f *FilterWrap) Bool(s string) {
	f.boolClause = s
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
	TermsMap        map[string]interface{} `json:"terms,omitempty"`
	TermMap         map[string]interface{} `json:"term,omitempty"`
	RangeMap        map[string]RangeFilter `json:"range,omitempty"`
	ExistsProp      *PropertyPathMarker    `json:"exists,omitempty"`
	MissingProp     *PropertyPathMarker    `json:"missing,omitempty"`
	AndFilters      []*FilterOp            `json:"and,omitempty"`
	OrFilters       []*FilterOp            `json:"or,omitempty"`
	NotFilters      []*FilterOp            `json:"not,omitempty"`
	LimitProp       *LimitFilter           `json:"limit,omitempty"`
	TypeProp        *TypeFilter            `json:"type,omitempty"`
	IdsProp         *IdsFilter             `json:"ids,omitempty"`
	ScriptProp      *ScriptFilter          `json:"script,omitempty"`
	GeoDistMap      map[string]interface{} `json:"geo_distance,omitempty"`
	GeoDistRangeMap map[string]interface{} `json:"geo_distance_range,omitempty"`
}

type PropertyPathMarker struct {
	Field string `json:"field"`
}

type LimitFilter struct {
	Value int `json:"value"`
}

type TypeFilter struct {
	Value string `json:"value"`
}

type IdsFilter struct {
	Type   []string      `json:"type,omitempty"`
	Values []interface{} `json:"values,omitempty"`
}

type ScriptFilter struct {
	Script   string                 `json:"script"`
	Params   map[string]interface{} `json:"params,omitempty"`
	IsCached bool                   `json:"_cache,omitempty"`
}

type RangeFilter struct {
	Gte      interface{} `json:"gte,omitempty"`
	Lte      interface{} `json:"lte,omitempty"`
	Gt       interface{} `json:"gt,omitempty"`
	Lt       interface{} `json:"lt,omitempty"`
	TimeZone string      `json:"time_zone,omitempty"` //Ideally this would be an int
}

type GeoLocation struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

type GeoField struct {
	GeoLocation
	Field string
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

func (f *FilterOp) And(filters ...*FilterOp) *FilterOp {
	if len(f.AndFilters) == 0 {
		f.AndFilters = filters[:]
	} else {
		f.AndFilters = append(f.AndFilters, filters...)
	}

	return f
}

func (f *FilterOp) Or(filters ...*FilterOp) *FilterOp {
	if len(f.OrFilters) == 0 {
		f.OrFilters = filters[:]
	} else {
		f.OrFilters = append(f.OrFilters, filters...)
	}

	return f
}

func (f *FilterOp) Not(filters ...*FilterOp) *FilterOp {
	if len(f.NotFilters) == 0 {
		f.NotFilters = filters[:]

	} else {
		f.NotFilters = append(f.NotFilters, filters...)
	}

	return f
}

func (f *FilterOp) GeoDistance(distance string, fields ...GeoField) *FilterOp {
	f.GeoDistMap = make(map[string]interface{})
	f.GeoDistMap["distance"] = distance
	for _, val := range fields {
		f.GeoDistMap[val.Field] = val.GeoLocation
	}

	return f
}

func (f *FilterOp) GeoDistanceRange(from string, to string, fields ...GeoField) *FilterOp {
	f.GeoDistRangeMap = make(map[string]interface{})
	f.GeoDistRangeMap["from"] = from
	f.GeoDistRangeMap["to"] = to

	for _, val := range fields {
		f.GeoDistRangeMap[val.Field] = val.GeoLocation
	}

	return f
}

// Helper to create values for the GeoDistance filters
func NewGeoField(field string, latitude float32, longitude float32) GeoField {
	return GeoField{
		GeoLocation: GeoLocation{Latitude: latitude, Longitude: longitude},
		Field:       field}
}

// Filter Terms
//
//   Filter().Terms("user","kimchy","stuff")
//	 Note: you can only have one terms clause in a filter. Use a bool filter to combine
func (f *FilterOp) Terms(field string, executionMode TermExecutionMode, values ...interface{}) *FilterOp {
	//You can only have one terms in a filter
	f.TermsMap = make(map[string]interface{})

	if executionMode != "" {
		f.TermsMap["execution"] = executionMode
	}

	f.TermsMap[field] = values

	return f
}

// Range adds a range filter for the given field.
func (f *FilterOp) Range(field string, gte interface{},
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

func (f *FilterOp) Type(fieldType string) *FilterOp {
	f.TypeProp = &TypeFilter{Value: fieldType}
	return f
}

func (f *FilterOp) Ids(ids ...interface{}) *FilterOp {
	f.IdsProp = &IdsFilter{Values: ids}
	return f
}

func (f *FilterOp) IdsByTypes(types []string, ids ...interface{}) *FilterOp {
	f.IdsProp = &IdsFilter{Type: types, Values: ids}
	return f
}

func (f *FilterOp) Exists(field string) *FilterOp {
	f.ExistsProp = &PropertyPathMarker{Field: field}
	return f
}

func (f *FilterOp) Missing(field string) *FilterOp {
	f.MissingProp = &PropertyPathMarker{Field: field}
	return f
}

func (f *FilterOp) Limit(maxResults int) *FilterOp {
	f.LimitProp = &LimitFilter{Value: maxResults}
	return f
}
