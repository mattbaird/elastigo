package core

import (
	u "github.com/araddon/gou"
	"testing"
)

func TestSearchRequest(t *testing.T) {
	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]string{"actor": "a*"},
		},
	}
	out, err := SearchRequest(true, "github", "", qry, "",0)
	//log.Println(out)
	Assert(&out != nil && err == nil, t, "Should get docs")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs but was %v", out.Hits.Len())
	Assert(u.CloseInt(out.Hits.Total, 588), t, "Should have 588 hits but was %v", out.Hits.Total)
}
