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

package search

import (
	"fmt"
	"github.com/araddon/gou"
	"github.com/bmizerany/assert"
	"github.com/mattbaird/elastigo/core"
	"log"
	"testing"
)

var (
	_ = log.Ldate
)

func TestSearchRequest(t *testing.T) {
	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]string{"actor": "a*"},
		},
	}
	out, err := core.SearchRequest("github", "", nil, qry)
	//log.Println(out)
	assert.T(t, &out != nil && err == nil, t, "Should get docs")
	assert.T(t, out.Hits.Total == 616 && out.Hits.Len() == 10, t, "Should have 616 hits but was %v", out.Hits.Total)
}

func TestSearchSimple(t *testing.T) {

	// searching without faceting
	qry := Search("github").Pretty().Query(
		Query().Search("add"),
	)
	out, _ := qry.Result()
	// how many different docs used the word "add"
	assert.T(t, out.Hits.Len() == 10, fmt.Sprintf("Should have 10 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 494, fmt.Sprintf("Should have 494 total= %v", out.Hits.Total))

	// now the same result from a "Simple" search
	out, _ = Search("github").Search("add").Result()
	assert.T(t, out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	assert.T(t, out.Hits.Total == 494, t, "Should have 494 total= %v", out.Hits.Total)
}

func TestSearchRequestQueryString(t *testing.T) {
	out, err := core.SearchUri("github", "", map[string]interface{}{"q": "actor:a*"})
	//log.Println(out)
	assert.T(t, &out != nil && err == nil, "Should get docs")
	assert.T(t, out.Hits.Total == 616, fmt.Sprintf("Should have 616 hits but was %v", out.Hits.Total))
}

func TestSearchFacetOne(t *testing.T) {
	/*
		A faceted search for what "type" of events there are
		- since we are not specifying an elasticsearch type it searches all ()

		{
		    "terms" : {
		      "_type" : "terms",
		      "missing" : 0,
		      "total" : 7561,
		      "other" : 0,
		      "terms" : [ {
		        "term" : "pushevent",
		        "count" : 4185
		      }, {
		        "term" : "createevent",
		        "count" : 786
		      }.....]
		    }
		 }

	*/
	qry := Search("github").Pretty().Facet(
		Facet().Fields("type").Size("25"),
	).Query(
		Query().All(),
	).Size("1")
	out, err := qry.Result()
	//log.Println(string(out.Facets))
	//gou.Debug(out)
	assert.T(t, out != nil && err == nil, "Should have output")
	if out == nil {
		t.Fail()
		return
	}
	h := gou.NewJsonHelper(out.Facets)
	assert.T(t, h.Int("type.total") == 8084, fmt.Sprintf("Should have 8084 results %v", h.Int("type.total")))
	assert.T(t, len(h.List("type.terms")) == 16, fmt.Sprintf("Should have 16 event types, %v", len(h.List("type.terms"))))

	// Now, lets try changing size to 10
	qry.FacetVal.Size("10")
	out, err = qry.Result()
	h = gou.NewJsonHelper(out.Facets)

	// still same doc count
	assert.T(t, h.Int("type.total") == 8084, fmt.Sprintf("Should have 8084 results %v", h.Int("type.total")))
	// make sure size worked
	assert.T(t, len(h.List("type.terms")) == 10, fmt.Sprintf("Should have 10 event types, %v", len(h.List("type.terms"))))

	// now, lets add a type (out of the 16)
	out, _ = Search("github").Type("IssueCommentEvent").Pretty().Facet(
		Facet().Fields("type").Size("25"),
	).Query(
		Query().All(),
	).Result()
	h = gou.NewJsonHelper(out.Facets)
	//log.Println(string(out.Facets))
	// still same doc count
	assert.T(t, h.Int("type.total") == 685, fmt.Sprintf("Should have 685 results %v", h.Int("type.total")))
	// we should only have one facettype because we limited to one type
	assert.T(t, len(h.List("type.terms")) == 1, fmt.Sprintf("Should have 1 event types, %v", len(h.List("type.terms"))))

	// now, add a second type (chained)
	out, _ = Search("github").Type("IssueCommentEvent").Type("PushEvent").Pretty().Facet(
		Facet().Fields("type").Size("25"),
	).Query(
		Query().All(),
	).Result()
	h = gou.NewJsonHelper(out.Facets)
	//log.Println(string(out.Facets))
	// still same doc count
	assert.T(t, h.Int("type.total") == 4941, fmt.Sprintf("Should have 4941 results %v", h.Int("type.total")))
	// make sure we now have 2 types
	assert.T(t, len(h.List("type.terms")) == 2, fmt.Sprintf("Should have 2 event types, %v", len(h.List("type.terms"))))

	//and instead of faceting on type, facet on userid
	// now, add a second type (chained)
	out, _ = Search("github").Type("IssueCommentEvent,PushEvent").Pretty().Facet(
		Facet().Fields("actor").Size("500"),
	).Query(
		Query().All(),
	).Result()
	h = gou.NewJsonHelper(out.Facets)
	// still same doc count
	assert.T(t, h.Int("actor.total") == 5158, t, "Should have 5158 results %v", h.Int("actor.total"))
	// make sure size worked
	assert.T(t, len(h.List("actor.terms")) == 500, t, "Should have 500 users, %v", len(h.List("actor.terms")))

}

func TestSearchFacetRange(t *testing.T) {
	// ok, now lets try facet but on actor field with a range
	qry := Search("github").Pretty().Facet(
		Facet().Fields("actor").Size("500"),
	).Query(
		Query().Search("add"),
	)
	out, err := qry.Result()
	assert.T(t, out != nil && err == nil, t, "Should have output")

	if out == nil {
		t.Fail()
		return
	}
	//log.Println(string(out.Facets))
	h := gou.NewJsonHelper(out.Facets)
	// how many different docs used the word "add", during entire time range
	assert.T(t, h.Int("actor.total") == 521, fmt.Sprintf("Should have 521 results %v", h.Int("actor.total")))
	// make sure size worked
	assert.T(t, len(h.List("actor.terms")) == 366,
		fmt.Sprintf("Should have 366 unique userids, %v", len(h.List("actor.terms"))))

	// ok, repeat but with a range showing different results
	qry = Search("github").Pretty().Facet(
		Facet().Fields("actor").Size("500"),
	).Query(
		Query().Range(
			Range().Field("created_at").From("2012-12-10T15:00:00-08:00").To("2012-12-10T15:10:00-08:00"),
		).Search("add"),
	)
	out, err = qry.Result()
	assert.T(t, out != nil && err == nil, t, "Should have output")

	if out == nil {
		t.Fail()
		return
	}
	//log.Println(string(out.Facets))
	h = gou.NewJsonHelper(out.Facets)
	// how many different events used the word "add", during time range?
	assert.T(t, h.Int("actor.total") == 97, fmt.Sprintf("Should have 97 results %v", h.Int("actor.total")))
	// make sure size worked
	assert.T(t, len(h.List("actor.terms")) == 71,
		fmt.Sprintf("Should have 71 event types, %v", len(h.List("actor.terms"))))

}

func TestSearchTerm(t *testing.T) {

	// ok, now lets try searching with term query (specific field/term)
	qry := Search("github").Query(
		Query().Term("repository.name", "jasmine"),
	)
	out, _ := qry.Result()
	// how many different docs have jasmine in repository.name?
	assert.T(t, out.Hits.Len() == 4, fmt.Sprintf("Should have 4 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 4, fmt.Sprintf("Should have 4 total= %v", out.Hits.Total))

}

func TestSearchFields(t *testing.T) {
	// same as terms, search using fields:
	//    how many different docs have jasmine in repository.name?
	qry := Search("github").Query(
		Query().Fields("repository.name", "jasmine", "", ""),
	)
	out, _ := qry.Result()

	assert.T(t, out.Hits.Len() == 4, fmt.Sprintf("Should have 4 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 4, fmt.Sprintf("Should have 4 total= %v", out.Hits.Total))
}

func TestSearchMissingExists(t *testing.T) {
	// search for docs that are missing repository.name
	qry := Search("github").Filter(
		Filter().Exists("repository.name"),
	)
	out, _ := qry.Result()
	assert.T(t, out.Hits.Len() == 10, fmt.Sprintf("Should have 10 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 7695, fmt.Sprintf("Should have 7695 total= %v", out.Hits.Total))

	qry = Search("github").Filter(
		Filter().Missing("repository.name"),
	)
	out, _ = qry.Result()
	assert.T(t, out.Hits.Len() == 10, fmt.Sprintf("Should have 10 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 389, fmt.Sprintf("Should have 389 total= %v", out.Hits.Total))
}

func TestSearchFilterQuery(t *testing.T) {

	// compound query + filter with query being wildcard
	out, _ := Search("github").Size("25").Query(
		Query().Fields("repository.name", "jas*", "", ""),
	).Filter(
		Filter().Terms("repository.has_wiki", true),
	).Result()
	if out == nil || &out.Hits == nil {
		t.Fail()
		return
	}

	assert.T(t, out.Hits.Len() == 7, fmt.Sprintf("Should have 7 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 7, fmt.Sprintf("Should have total=7 but was %v", out.Hits.Total))
}

func TestSearchRange(t *testing.T) {

	// now lets filter by a subset of the total time
	out, _ := Search("github").Size("25").Query(
		Query().Range(
			Range().Field("created_at").From("2012-12-10T15:00:00-08:00").To("2012-12-10T15:10:00-08:00"),
		).Search("add"),
	).Result()
	assert.T(t, out != nil && &out.Hits != nil, "Must not have nil results, or hits")
	assert.T(t, out.Hits.Len() == 25, fmt.Sprintf("Should have 25 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 92, fmt.Sprintf("Should have total=92 but was %v", out.Hits.Total))
}

func TestSearchSortOrder(t *testing.T) {

	// ok, now lets try sorting by repository watchers descending
	qry := Search("github").Pretty().Query(
		Query().All(),
	).Sort(
		Sort("repository.watchers").Desc(),
	)
	out, _ := qry.Result()

	// how many different docs used the word "add", during entire time range
	assert.T(t, out.Hits.Len() == 10, fmt.Sprintf("Should have 10 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 8084, fmt.Sprintf("Should have 8084 total= %v", out.Hits.Total))
	h1 := gou.NewJsonHelper(out.Hits.Hits[0].Source)
	assert.T(t, h1.Int("repository.watchers") == 41377,
		fmt.Sprintf("Should have 41377 watchers= %v", h1.Int("repository.watchers")))

	// ascending
	out, _ = Search("github").Pretty().Query(
		Query().All(),
	).Sort(
		Sort("repository.watchers"),
	).Result()
	// how many different docs used the word "add", during entire time range
	assert.T(t, out.Hits.Len() == 10, fmt.Sprintf("Should have 10 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 8084, fmt.Sprintf("Should have 8084 total= %v", out.Hits.Total))
	h2 := gou.NewJsonHelper(out.Hits.Hits[0].Source)
	assert.T(t, h2.Int("repository.watchers") == 0,
		fmt.Sprintf("Should have 0 watchers= %v", h2.Int("repository.watchers")))

	// sort descending with search
	out, _ = Search("github").Pretty().Size("5").Query(
		Query().Search("python"),
	).Sort(
		Sort("repository.watchers").Desc(),
	).Result()
	//log.Println(out)
	//log.Println(err)
	// how many different docs used the word "add", during entire time range
	assert.T(t, out.Hits.Len() == 5, fmt.Sprintf("Should have 5 docs %v", out.Hits.Len()))
	assert.T(t, out.Hits.Total == 734, fmt.Sprintf("Should have 734 total= %v", out.Hits.Total))
	h3 := gou.NewJsonHelper(out.Hits.Hits[0].Source)
	assert.T(t, h3.Int("repository.watchers") == 8659,
		fmt.Sprintf("Should have 8659 watchers= %v", h3.Int("repository.watchers")))

}
