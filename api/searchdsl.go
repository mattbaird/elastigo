package api

type SearchRequest struct {
	From   int    `json:"from,omitempty"`
	Size   int    `json:"size,omitempty"`
	Query  Query  `json:"query,omitempty"`
	Filter Filter `json:"filter,omitempty"`
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
