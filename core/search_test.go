package core

import (
	"testing"
)

func TestSearchRequest(t *testing.T) {
	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]string{"actor": "a*"},
		},
	}
	out, err := SearchRequest(true, "github", "", qry, "")
	//log.Println(out)
	Assert(&out != nil && err == nil, t, "Should get docs")
	Assert(out.Hits.Total == 588 && out.Hits.Len() == 10, t, "Should have 588 hits but was %v", out.Hits.Total)
}
