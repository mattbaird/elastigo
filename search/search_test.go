package search

import (
	"flag"
	"fmt"
	"github.com/araddon/gou"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"log"
	"testing"
)

var (
	_                 = log.Ldate
	hasStartedTesting bool
	eshost            *string = flag.String("host", "localhost", "Elasticsearch Server Host Address")
)

// dumb simple assert for testing, printing
//    Assert(len(items) == 9, t, "Should be 9 but was %d", len(items))
func Assert(is bool, t *testing.T, format string, args ...interface{}) {
	if is == false {
		log.Printf(format, args...)
		t.Fail()
	}
}

/*

usage:

	test -v -host eshost 

*/

func init() {
	InitTests(false)
	core.DebugRequests = true
}

func InitTests(startIndexor bool) {
	if !hasStartedTesting {
		flag.Parse()
		hasStartedTesting = true
		flag.Parse()
		log.SetFlags(log.Ltime | log.Lshortfile)
		api.Domain = *eshost
	}
}

func TestSearchRequest(t *testing.T) {
	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]string{"actor": "a*"},
		},
	}
	out, err := core.SearchRequest(true, "github", "", qry, "")
	//log.Println(out)
	Assert(&out != nil && err == nil, t, "Should get docs")
	Assert(out.Hits.Total == 588 && out.Hits.Len() == 10, t, "Should have 588 hits but was %v", out.Hits.Total)
}

func TestSearchSimple(t *testing.T) {

	// searching without faceting
	qry := Search("github").Pretty().Query(
		Query().Search("add"),
	)
	out, _ := qry.Result()
	// how many different docs used the word "add" 
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 478, t, "Should have 478 total= %v", out.Hits.Total)

	// now the same result from a "Simple" search
	out, _ = Search("github").Search("add").Result()
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 478, t, "Should have 478 total= %v", out.Hits.Total)
}

func TestSearchRequestQueryString(t *testing.T) {
	out, err := core.SearchUri("github", "", "actor:a*", "")
	//log.Println(out)
	Assert(&out != nil && err == nil, t, "Should get docs")
	Assert(out.Hits.Total == 588, t, "Should have 588 hits but was %v", out.Hits.Total)
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
	)
	out, err := qry.Result()
	//log.Println(string(out.Facets))
	Assert(out != nil && err == nil, t, "Should have output")
	if out == nil {
		t.Fail()
		return
	}
	h := gou.NewJsonHelper(out.Facets)
	Assert(h.Int("terms.total") == 7545, t, "Should have 7545 results %v", h.Int("terms.total"))
	Assert(len(h.List("terms.terms")) == 16, t, "Should have 16 event types, %v", len(h.List("terms.terms")))

	// Now, lets try changing size to 10
	qry.FacetVal.Size("10")
	out, err = qry.Result()
	h = gou.NewJsonHelper(out.Facets)

	// still same doc count
	Assert(h.Int("terms.total") == 7545, t, "Should have 7545 results %v", h.Int("terms.total"))
	// make sure size worked
	Assert(len(h.List("terms.terms")) == 10, t, "Should have 10 event types, %v", len(h.List("terms.terms")))

	// now, lets add a type (out of the 16) 
	out, _ = Search("github").Type("IssueCommentEvent").Pretty().Facet(
		Facet().Fields("type").Size("25"),
	).Query(
		Query().All(),
	).Result()
	h = gou.NewJsonHelper(out.Facets)
	//log.Println(string(out.Facets))
	// still same doc count
	Assert(h.Int("terms.total") == 679, t, "Should have 679 results %v", h.Int("terms.total"))
	// we should only have one facettype because we limited to one type
	Assert(len(h.List("terms.terms")) == 1, t, "Should have 1 event types, %v", len(h.List("terms.terms")))

	// now, add a second type (chained)
	out, _ = Search("github").Type("IssueCommentEvent").Type("PushEvent").Pretty().Facet(
		Facet().Fields("type").Size("25"),
	).Query(
		Query().All(),
	).Result()
	h = gou.NewJsonHelper(out.Facets)
	//log.Println(string(out.Facets))
	// still same doc count
	Assert(h.Int("terms.total") == 4863, t, "Should have 4863 results %v", h.Int("terms.total"))
	// make sure we now have 2 types
	Assert(len(h.List("terms.terms")) == 2, t, "Should have 2 event types, %v", len(h.List("terms.terms")))

	//and instead of faceting on type, facet on userid
	// now, add a second type (chained)
	out, _ = Search("github").Type("IssueCommentEvent,PushEvent").Pretty().Facet(
		Facet().Fields("actor").Size("500"),
	).Query(
		Query().All(),
	).Result()
	h = gou.NewJsonHelper(out.Facets)
	// still same doc count
	Assert(h.Int("terms.total") == 5065, t, "Should have 5065 results %v", h.Int("terms.total"))
	// make sure size worked
	Assert(len(h.List("terms.terms")) == 500, t, "Should have 500 users, %v", len(h.List("terms.terms")))

}

func TestSearchFacetRange(t *testing.T) {
	// ok, now lets try facet but on actor field with a range
	qry := Search("github").Pretty().Facet(
		Facet().Fields("actor").Size("500"),
	).Query(
		Query().Search("add"),
	)
	out, err := qry.Result()
	Assert(out != nil && err == nil, t, "Should have output")

	if out == nil {
		t.Fail()
		return
	}
	//log.Println(string(out.Facets))
	h := gou.NewJsonHelper(out.Facets)
	// how many different docs used the word "add", during entire time range
	Assert(h.Int("terms.total") == 504, t, "Should have 504 results %v", h.Int("terms.total"))
	// make sure size worked
	Assert(len(h.List("terms.terms")) == 361, t, "Should have 361 unique userids, %v", len(h.List("terms.terms")))

	// ok, repeat but with a range showing different results
	qry = Search("github").Pretty().Facet(
		Facet().Fields("actor").Size("500"),
	).Query(
		Query().Range(
			Range().Field("created_at").From("2012-12-10T15:00:00-08:00").To("2012-12-10T15:10:00-08:00"),
		).Search("add"),
	)
	out, err = qry.Result()
	Assert(out != nil && err == nil, t, "Should have output")

	if out == nil {
		t.Fail()
		return
	}
	//log.Println(string(out.Facets))
	h = gou.NewJsonHelper(out.Facets)
	// how many different events used the word "add", during time range?
	Assert(h.Int("terms.total") == 97, t, "Should have 97 results %v", h.Int("terms.total"))
	// make sure size worked
	Assert(len(h.List("terms.terms")) == 71, t, "Should have 71 event types, %v", len(h.List("terms.terms")))

}

func TestSearchTerm(t *testing.T) {

	// ok, now lets try searching with term query (specific field/term)
	qry := Search("github").Query(
		Query().Term("repository.name", "jasmine"),
	)
	out, _ := qry.Result()
	// how many different docs have jasmine in repository.name?
	Assert(out.Hits.Len() == 3, t, "Should have 3 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 3, t, "Should have 3 total= %v", out.Hits.Total)

}

func TestSearchFields(t *testing.T) {
	// same as terms, search using fields:
	//    how many different docs have jasmine in repository.name?
	qry := Search("github").Query(
		Query().Fields("repository.name", "jasmine", "", ""),
	)
	out, _ := qry.Result()

	Assert(out.Hits.Len() == 3, t, "Should have 3 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 3, t, "Should have 3 total= %v", out.Hits.Total)
}

func TestSearchMissingExists(t *testing.T) {
	// search for docs that are missing repository.name
	qry := Search("github").Filter(
		Filter().Exists("repository.name"),
	)
	out, _ := qry.Result()
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 7241, t, "Should have 7241 total= %v", out.Hits.Total)

	qry = Search("github").Filter(
		Filter().Missing("repository.name"),
	)
	out, _ = qry.Result()
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 304, t, "Should have 304 total= %v", out.Hits.Total)
}

func TestSearchRange(t *testing.T) {

	// now lets filter by a subset of the total time
	out, _ := Search("github").Size("25").Query(
		Query().Range(
			Range().Field("created_at").From("2012-12-10T15:00:00-08:00").To("2012-12-10T15:10:00-08:00"),
		).Search("add"),
	).Result()
	fmt.Println(out)
	if out == nil || &out.Hits == nil {
		t.Fail()
		return
	}

	Assert(out.Hits.Len() == 25, t, "Should have 25 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 92, t, "Should have total=92 but was %v", out.Hits.Total)
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
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 7545, t, "Should have 7545 total= %v", out.Hits.Total)
	h1 := gou.NewJsonHelper(out.Hits.Hits[0].Source)
	Assert(h1.Int("repository.watchers") == 41377, t, "Should have 41377 watchers= %v", h1.Int("repository.watchers"))

	// ascending 
	out, _ = Search("github").Pretty().Query(
		Query().All(),
	).Sort(
		Sort("repository.watchers"),
	).Result()
	// how many different docs used the word "add", during entire time range
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 7545, t, "Should have 7545 total= %v", out.Hits.Total)
	h2 := gou.NewJsonHelper(out.Hits.Hits[0].Source)
	Assert(h2.Int("repository.watchers") == 0, t, "Should have 0 watchers= %v", h2.Int("repository.watchers"))

	// sort descending with search 
	out, _ = Search("github").Pretty().Size("5").Query(
		Query().Search("python"),
	).Sort(
		Sort("repository.watchers").Desc(),
	).Result()
	//log.Println(out)
	//log.Println(err)
	// how many different docs used the word "add", during entire time range
	Assert(out.Hits.Len() == 5, t, "Should have 5 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 708, t, "Should have 708 total= %v", out.Hits.Total)
	h3 := gou.NewJsonHelper(out.Hits.Hits[0].Source)
	Assert(h3.Int("repository.watchers") == 8659, t, "Should have 8659 watchers= %v", h3.Int("repository.watchers"))

}
