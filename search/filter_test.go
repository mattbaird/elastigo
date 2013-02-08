package search

import (
	//"encoding/json"
	. "github.com/araddon/gou"
	"testing"
)

func TestFilters(t *testing.T) {
	// search for docs that are missing repository.name
	qry := Search("github").Filter(
		Filter().Exists("repository.name"),
	)
	out, err := qry.Result()
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 7241, t, "Should have 7241 total= %v", out.Hits.Total)

	qry = Search("github").Filter(
		Filter().Missing("repository.name"),
	)
	out, _ = qry.Result()
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 304, t, "Should have 304 total= %v", out.Hits.Total)

	//actor_attributes: {type: "User",
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
	)
	out, _ = qry.Result()
	Debug(out)
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 65, t, "Should have 65 total= %v", out.Hits.Total)

	/*
		Should this be an AND by default?
	*/
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
		Filter().Terms("repository.has_wiki", true),
	)
	out, err = qry.Result()
	Debug(out)
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 43, t, "Should have 43 total= %v", out.Hits.Total)

	// NOW, lets try with two query calls instead of one
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
	)
	qry.Filter(
		Filter().Terms("repository.has_wiki", true),
	)
	out, err = qry.Result()
	Debug(out)
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 43, t, "Should have 43 total= %v", out.Hits.Total)

	qry = Search("github").Filter(
		"or",
		Filter().Terms("actor_attributes.location", "portland"),
		Filter().Terms("repository.has_wiki", true),
	)
	out, err = qry.Result()
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 6290, t, "Should have 6290 total= %v", out.Hits.Total)
}
