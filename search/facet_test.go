package search

import (
	//"encoding/json"
	. "github.com/araddon/gou"
	"testing"
)

func TestFacetRegex(t *testing.T) {

	// This is a possible solution for auto-complete
	out, _ := Search("github").Size("0").Facet(
		Facet().Regex("repository.name", "no.*").Size("8"),
	).Result()
	if out == nil || &out.Hits == nil {
		t.Fail()
		return
	}
	//Debug(string(out.Facets))
	fh := NewJsonHelper([]byte(out.Facets))
	facets := fh.Helpers("/repository.name/terms")
	Assert(len(facets) == 8, t, "Should have 8? but was %v", len(facets))
	// for _, f := range facets {
	// 	Debug(f)
	// }
}
