package elastigo

func NewHighlight() *HighlightDsl {
	return &HighlightDsl{}
}

type HighlightDsl struct {
	Options   *HighlightOptions
	TagSchema string                      `json:"tag_schema,omitempty"`
	Fields    map[string]HighlightOptions `json:"fields,omitempty"`
}

func NewHighlightOpts() *HighlightOptions {
	return &HighlightOptions{}
}

type HighlightOptions struct {
	BoundaryCharsVal   string    `json:"boundary_chars,omitempty"`
	BoundaryMaxScanVal int       `json:"boundary_max_scan,omitempty"`
	PreTags            []string  `json:"pre_tags,omitempty"`
	PostTags           []string  `json:"post_tags,omitempty"`
	FragmentSizeVal    int       `json:"fragment_size,omitempty"`
	NumOfFragmentsVal  int       `json:"number_of_fragments,omitempty"`
	HighlightQuery     *QueryDsl `json:"highlight_query,omitempty"`
	MatchedFieldsVal   []string  `json:"matched_fields,omitempty"`
	OrderVal           string    `json:"order,omitempty"`
	TypeVal            string    `json:"type,omitempty"`
}

func (h *HighlightDsl) AddField(name string, settings *HighlightOptions) *HighlightDsl {
	if h.Fields == nil {
		h.Fields = make(map[string]HighlightOptions)
	}

	if settings != nil {
		h.Fields[name] = *settings
	} else {
		h.Fields[name] = HighlightOptions{}
	}

	return h
}



func (h *HighlightDsl) Schema(schema string) *HighlightDsl {
	h.TagSchema = schema
	return h
}

func (h *HighlightDsl) Settings(options *HighlightOptions) *HighlightDsl {
	h.Options = options
	return h
}

func (o *HighlightOptions) BoundaryChars(chars string) *HighlightOptions {
	o.BoundaryCharsVal = chars
	return o
}

func (o *HighlightOptions) Tags(pre string, post string) *HighlightOptions {
	if o.PreTags == nil {
		o.PreTags = []string{pre}
	} else {
		o.PreTags = append(o.PreTags, pre)
	}

	if o.PostTags == nil {
		o.PostTags = []string{post}
	} else {
		o.PostTags = append(o.PostTags, post)
	}

	return o
}
