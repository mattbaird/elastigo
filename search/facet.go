package search

import (
	"encoding/json"

	. "github.com/araddon/gou"
)

var (
	_ = DEBUG
)

/*
"facets": {
    "terms": {
		"terms": {
			"field": [
			  "@fields.category"
			],
			"size": 25
		}
    }
}


"facets": {
  "actors": { "terms": {"field": ["actor"],"size": "10" }}
  , "langauge": { "terms": {"field": ["repository.language"],"size": "10" }}
}

*/
func Facet() *FacetDsl {
	return &FacetDsl{}
}

type FacetDsl struct {
	size  string
	Terms map[string]Term `json:"terms,omitempty"`
}

func (m *FacetDsl) Size(size string) *FacetDsl {
	m.size = size
	return m
}

func (m *FacetDsl) Regex(field, match string) *FacetDsl {
	if len(m.Terms) == 0 {
		m.Terms = make(map[string]Term)
	}
	m.Terms[field] = Term{Terms{Fields: []string{field}, Regex: match}}
	return m
}

func (m *FacetDsl) Fields(fields ...string) *FacetDsl {
	if len(fields) < 1 {
		return m
	}
	if len(m.Terms) == 0 {
		m.Terms = make(map[string]Term)
	}
	m.Terms[fields[0]] = Term{Terms{Fields: fields}}
	return m
}

func (m *FacetDsl) MarshalJSON() ([]byte, error) {
	// Custom marshall
	for _, t := range m.Terms {
		t.Terms.Size = m.size
	}
	return json.Marshal(&m.Terms)
}
