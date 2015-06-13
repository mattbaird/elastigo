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
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

func TestAndDsl(t *testing.T) {
	filter := Filter().And(Filter().Term("test", "asdf")).
		And(Filter().Range("rangefield", 1, 2, 3, 4, "+08:00"))
	actual := GetJson(filter)

	actualFilters := actual["and"].([]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, 2, len(actualFilters), "Should have 2 filters")
	assert.Equal(t, true, HasKey(actualFilters[0].(map[string]interface{}), "term"), "first filter is term")
	assert.Equal(t, true, HasKey(actualFilters[1].(map[string]interface{}), "range"), "second filter is range")
}

func TestOrDsl(t *testing.T) {
	filter := Filter().Or(Filter().Term("test", "asdf")).
		Or(Filter().Range("rangefield", 1, 2, 3, 4, "+08:00"))
	actual := GetJson(filter)

	actualFilters := actual["or"].([]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, 2, len(actualFilters), "Should have 2 filters")
	assert.Equal(t, true, HasKey(actualFilters[0].(map[string]interface{}), "term"), "first filter is term")
	assert.Equal(t, true, HasKey(actualFilters[1].(map[string]interface{}), "range"), "second filter is range")
}

func TestNotDsl(t *testing.T) {
	filter := Filter().Not(Filter().Term("test", "asdf")).
		Not(Filter().Range("rangefield", 1, 2, 3, 4, "+08:00"))
	actual := GetJson(filter)

	actualFilters := actual["not"].([]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, 2, len(actualFilters), "Should have 2 filters")
	assert.Equal(t, true, HasKey(actualFilters[0].(map[string]interface{}), "term"), "first filter is term")
	assert.Equal(t, true, HasKey(actualFilters[1].(map[string]interface{}), "range"), "second filter is range")
}

func TestTermsDsl(t *testing.T) {
	filter := Filter().Terms("Sample", TEM_AND, "asdf", 123, true)
	actual := GetJson(filter)

	actualTerms := actual["terms"].(map[string]interface{})
	actualValues := actualTerms["Sample"].([]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, 3, len(actualValues), "Should have 3 term values")
	assert.Equal(t, actualValues[0], "asdf")
	assert.Equal(t, actualValues[1], float64(123))
	assert.Equal(t, actualValues[2], true)
	assert.Equal(t, "and", actualTerms["execution"])
}

func TestTermDsl(t *testing.T) {
	filter := Filter().Term("Sample", "asdf").Term("field2", 341.4)
	actual := GetJson(filter)

	actualTerm := actual["term"].(map[string]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, "asdf", actualTerm["Sample"])
	assert.Equal(t, float64(341.4), actualTerm["field2"])
}

func TestRangeDsl(t *testing.T) {
	filter := Filter().Range("rangefield", 1, 2, 3, 4, "+08:00")
	actual := GetJson(filter)
	//A bit lazy, probably should assert keys exist
	actualRange := actual["range"].(map[string]interface{})["rangefield"].(map[string]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, float64(1), actualRange["gte"])
	assert.Equal(t, float64(2), actualRange["gt"])
	assert.Equal(t, float64(3), actualRange["lte"])
	assert.Equal(t, float64(4), actualRange["lt"])
	assert.Equal(t, "+08:00", actualRange["time_zone"])
}

func TestExistsDsl(t *testing.T) {
	filter := Filter().Exists("field1")
	actual := GetJson(filter)

	actualValue := actual["exists"].(map[string]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, "field1", actualValue["field"], "exist field should match")
}

func TestMissingDsl(t *testing.T) {
	filter := Filter().Missing("field1")
	actual := GetJson(filter)

	actualValue := actual["missing"].(map[string]interface{})

	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, "field1", actualValue["field"], "missing field should match")
}

func TestLimitDsl(t *testing.T) {
	filter := Filter().Limit(100)
	actual := GetJson(filter)

	actualValue := actual["limit"].(map[string]interface{})
	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, float64(100), actualValue["value"], "limit value should match")
}

func TestTypeDsl(t *testing.T) {
	filter := Filter().Type("my_type")
	actual := GetJson(filter)

	actualValue := actual["type"].(map[string]interface{})
	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, "my_type", actualValue["value"], "type value should match")
}

func TestIdsDsl(t *testing.T) {
	filter := Filter().Ids("test", "asdf", "fdsa")
	actual := GetJson(filter)

	actualValue := actual["ids"].(map[string]interface{})
	actualValues := actualValue["values"].([]interface{})
	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, nil, actualValue["type"], "Should have no type specified")
	assert.Equal(t, 3, len(actualValues), "Should have 3 values specified")
	assert.Equal(t, "test", actualValues[0], "Should have same value")
	assert.Equal(t, "asdf", actualValues[1], "Should have same value")
	assert.Equal(t, "fdsa", actualValues[2], "Should have same value")
}

func TestIdsTypeDsl(t *testing.T) {
	filter := Filter().IdsByTypes([]string{"my_type"}, "test", "asdf", "fdsa")
	actual := GetJson(filter)

	actualValue := actual["ids"].(map[string]interface{})
	actualTypes := actualValue["type"].([]interface{})
	actualValues := actualValue["values"].([]interface{})
	assert.Equal(t, 1, len(actual), "JSON should only have one key")
	assert.Equal(t, 1, len(actualTypes), "Should have one type specified")
	assert.Equal(t, "my_type", actualTypes[0], "Should have correct type specified")
	assert.Equal(t, 3, len(actualValues), "Should have 3 values specified")
	assert.Equal(t, "test", actualValues[0], "Should have same value")
	assert.Equal(t, "asdf", actualValues[1], "Should have same value")
	assert.Equal(t, "fdsa", actualValues[2], "Should have same value")
}

func TestGeoDistDsl(t *testing.T) {
	filter := Filter().GeoDistance("100km", NewGeoField("pin.location", 32.3, 23.4))
	actual := GetJson(filter)

	actualValue := actual["geo_distance"].(map[string]interface{})
	actualLocation := actualValue["pin.location"].(map[string]interface{})
	assert.Equal(t, "100km", actualValue["distance"], "Distance should be equal")
	assert.Equal(t, float64(32.3), actualLocation["lat"], "Latitude should be equal")
	assert.Equal(t, float64(23.4), actualLocation["lon"], "Longitude should be equal")
}

func TestGeoDistRangeDsl(t *testing.T) {
	filter := Filter().GeoDistanceRange("100km", "200km", NewGeoField("pin.location", 32.3, 23.4))
	actual := GetJson(filter)

	actualValue := actual["geo_distance_range"].(map[string]interface{})
	actualLocation := actualValue["pin.location"].(map[string]interface{})
	assert.Equal(t, "100km", actualValue["from"], "From should be equal")
	assert.Equal(t, "200km", actualValue["to"], "To should be equal")
	assert.Equal(t, float64(32.3), actualLocation["lat"], "Latitude should be equal")
	assert.Equal(t, float64(23.4), actualLocation["lon"], "Longitude should be equal")
}

func TestFilters(t *testing.T) {
	c := NewTestConn()

	// search for docs that are missing repository.name
	qry := Search("github").Filter(
		Filter().Exists("repository.name"),
	)
	out, err := qry.Result(c)
	assert.T(t, err == nil, t, "should not have error")
	expectedDocs := 10
	expectedHits := 7695
	assert.T(t, out.Hits.Len() == expectedDocs, fmt.Sprintf("Should have %v docs got %v", expectedDocs, out.Hits.Len()))
	assert.T(t, out.Hits.Total == expectedHits, fmt.Sprintf("Should have %v total got %v", expectedHits, out.Hits.Total))
	qry = Search("github").Filter(
		Filter().Missing("repository.name"),
	)
	expectedHits = 390
	out, _ = qry.Result(c)
	assert.T(t, out.Hits.Len() == expectedDocs, fmt.Sprintf("Should have %v docs got %v", expectedDocs, out.Hits.Len()))
	assert.T(t, out.Hits.Total == expectedHits, fmt.Sprintf("Should have %v total got %v", expectedHits, out.Hits.Total))

	//actor_attributes: {type: "User",
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", TEM_DEFAULT, "portland"),
	)
	out, _ = qry.Result(c)
	expectedDocs = 10
	expectedHits = 71
	assert.T(t, out.Hits.Len() == expectedDocs, fmt.Sprintf("Should have %v docs got %v", expectedDocs, out.Hits.Len()))
	assert.T(t, out.Hits.Total == expectedHits, fmt.Sprintf("Should have %v total got %v", expectedHits, out.Hits.Total))

	/*
		Should this be an AND by default?
	*/
	qry = Search("github").Filter(
		Filter().And(Filter().Terms("actor_attributes.location", TEM_DEFAULT, "portland")).
			And(Filter().Terms("repository.has_wiki", TEM_DEFAULT, true)))
	out, err = qry.Result(c)
	expectedDocs = 10
	expectedHits = 44
	assert.T(t, err == nil, t, "should not have error")
	assert.T(t, out.Hits.Len() == expectedDocs, fmt.Sprintf("Should have %v docs got %v", expectedDocs, out.Hits.Len()))
	assert.T(t, out.Hits.Total == expectedHits, fmt.Sprintf("Should have %v total got %v", expectedHits, out.Hits.Total))

	// NOW, lets try with two query calls instead of one
	qry = Search("github").Filter(
		Filter().
			And(Filter().Terms("actor_attributes.location", TEM_DEFAULT, "portland")).
			And(Filter().Terms("repository.has_wiki", TEM_DEFAULT, true)),
	)

	out, err = qry.Result(c)
	//gou.Debug(out)
	assert.T(t, err == nil, t, "should not have error")
	assert.T(t, out.Hits.Len() == expectedDocs, fmt.Sprintf("Should have %v docs got %v", expectedDocs, out.Hits.Len()))
	assert.T(t, out.Hits.Total == expectedHits, fmt.Sprintf("Should have %v total got %v", expectedHits, out.Hits.Total))

	qry = Search("github").Filter(
		Filter().Or(Filter().Terms("actor_attributes.location", TEM_DEFAULT, "portland")).
			Or(Filter().Terms("repository.has_wiki", TEM_DEFAULT, true)),
	)
	out, err = qry.Result(c)
	expectedHits = 6676
	assert.T(t, err == nil, t, "should not have error")
	assert.T(t, out.Hits.Len() == expectedDocs, fmt.Sprintf("Should have %v docs got %v", expectedDocs, out.Hits.Len()))
	assert.T(t, out.Hits.Total == expectedHits, fmt.Sprintf("Should have %v total got %v", expectedHits, out.Hits.Total))
}

func TestFilterRange(t *testing.T) {
	c := NewTestConn()

	// now lets filter range for repositories with more than 100 forks
	out, _ := Search("github").Size("25").Filter(Filter().
		Range("repository.forks", 100, nil, nil, nil, "")).Result(c)

	if out == nil || &out.Hits == nil {
		t.Fail()
		return
	}
	expectedDocs := 25
	expectedHits := 725

	assert.T(t, out.Hits.Len() == expectedDocs, fmt.Sprintf("Should have %v docs got %v", expectedDocs, out.Hits.Len()))
	assert.T(t, out.Hits.Total == expectedHits, fmt.Sprintf("Should have total %v got %v", expectedHits, out.Hits.Total))
}
