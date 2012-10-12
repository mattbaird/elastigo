package api

type BaseResponse struct {
	Ok      bool        `json:"ok"`
	Index   string      `json:"_index"`
	Type    string      `json:"_type"`
	Id      string      `json:"_id"`
	Source  interface{} `json:"_source"` // depends on the schema you've defined
	Version int         `json:"_version"`
	Found   bool        `json:"found"`
	Exists  bool        `json:"exists"`
}

func Pretty(pretty bool) string {
	prettyString := ""
	if pretty == true {
		prettyString = "pretty=1"
	}
	return prettyString
}
