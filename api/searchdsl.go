package api

type SearchRequest struct {
	From   int    `json:"from,omitifempty"`
	Size   int    `json:"size,omitifempty"`
	Query  Query  `json:"query,omitifempty"`
	Filter Filter `json:"filter,omitifempty"`
}

type Filter struct {
	Term Term `json:"term"`
}

type Facets struct {
	Tag Terms `json:"tag"`
}

type Terms struct {
	Terms string `json:"terms"`
}
