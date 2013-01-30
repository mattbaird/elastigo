package search

import ()

/*
"filter": {
	"range": {
	  "@timestamp": {
	    "from": "2012-12-29T16:52:48+00:00",
	    "to": "2012-12-29T17:52:48+00:00"
	  }
	}
}
"filter": {
    "missing": {
        "field": "repository.name"
    }
}

*/

func Filter() *FilterOp {
	return &FilterOp{}
}

type FilterOp struct {
	curField    string
	Range       map[string]map[string]string `json:"range,omitempty"`
	Exist       map[string]string            `json:"exists,omitempty"`
	MisssingVal map[string]string            `json:"missing,omitempty"`
}

func Range() *FilterOp {
	return &FilterOp{Range: make(map[string]map[string]string)}
}

func (f *FilterOp) Field(fld string) *FilterOp {
	f.curField = fld
	if _, ok := f.Range[fld]; !ok {
		m := make(map[string]string)
		f.Range[fld] = m
	}
	return f
}
func (f *FilterOp) From(from string) *FilterOp {
	f.Range[f.curField]["from"] = from
	return f
}
func (f *FilterOp) To(to string) *FilterOp {
	f.Range[f.curField]["to"] = to
	return f
}
func (f *FilterOp) Exists(name string) *FilterOp {
	f.Exist = map[string]string{"field": name}
	return f
}
func (f *FilterOp) Missing(name string) *FilterOp {
	f.MisssingVal = map[string]string{"field": name}
	return f
}

// Add another Filterop, "combines" two filter ops into one
func (f *FilterOp) Add(fop *FilterOp) *FilterOp {
	if len(fop.Exist) > 0 {
		f.Exist = fop.Exist
	}
	if len(fop.MisssingVal) > 0 {
		f.MisssingVal = fop.MisssingVal
	}
	if len(fop.Range) > 0 {
		f.Range = fop.Range
	}
	return f
}
