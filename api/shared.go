package api

type Query struct {
	Query Term `json:"query"`
}

type Term struct {
	Term string `json:"term"`
}

func (q Query) setQuery(query string) {
	q.Query.Term = query
}

