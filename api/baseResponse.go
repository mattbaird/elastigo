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

func Pretty(pretty bool) string {
	prettyString := ""
	if pretty == true {
		prettyString = "pretty=1"
	}
	return prettyString
}
