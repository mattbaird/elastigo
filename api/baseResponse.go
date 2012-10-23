package api

type BaseResponse struct {
	Ok      bool        `json:"ok"`
	Index   string      `json:"_index,omitifempty"`
	Type    string      `json:"_type,omitifempty"`
	Id      string      `json:"_id,omitifempty"`
	Source  interface{} `json:"_source,omitifempty"` // depends on the schema you've defined
	Version int         `json:"_version,omitifempty"`
	Found   bool        `json:"found,omitifempty"`
	Exists  bool        `json:"exists,omitifempty"`
}

type Status struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

type Query struct {
	Query Term `json:"query"`
}

type Term struct {
	Term string `json:"term"`
}

func (q Query) setQuery(query string) {
	q.Query.Term = query
}

type Match struct {
	OK           bool         `json:"ok"`
	Matches      []string     `json:"matches"`
	Explaination Explaination `json:"explaination,omitifempty"`
}

type Explaination struct {
	Value       float32        `json:"value"`
	Description string         `json:"description"`
	Details     []Explaination `json:"details,omitifempty"`
}

func Pretty(pretty bool) string {
	prettyString := ""
	if pretty == true {
		prettyString = "pretty=1"
	}
	return prettyString
}
