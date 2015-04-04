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
	//"github.com/araddon/gou"
	"github.com/bmizerany/assert"
	"testing"
	"encoding/json"
)

func GetJson(input interface{}) map[string]interface{} {
	var result map[string]interface{}
	bytes, _ := json.Marshal(input)

	json.Unmarshal(bytes, &result)
	return result
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
