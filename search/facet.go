package search

import ()

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
*/
func Facet() *FacetDsl {
	return &FacetDsl{&Term{Terms{nil, ""}}}
}

type FacetDsl struct {
	TermsVal *Term `json:"terms,omitempty"`
}

func (f *FacetDsl) Size(size string) *FacetDsl {
	f.TermsVal.Terms.Size = size
	return f
}
func (f *FacetDsl) Fields(fields ...string) *FacetDsl {
	flds := make([]string, 0)
	for _, field := range fields {
		flds = append(flds, field)
	}
	f.TermsVal.Terms.Fields = flds
	return f
}
