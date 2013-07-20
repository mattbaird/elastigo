package core

import (
	"encoding/json"
	"fmt"
	"github.com/meanpath/elastigo/api"
)

// The more like this (mlt) API allows to get documents that are “like” a specified document. 
// http://www.elasticsearch.org/guide/reference/api/more-like-this.html
func MoreLikeThis(pretty bool, index string, _type string, id string, query MoreLikeThisQuery) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse
	url = fmt.Sprintf("/%s/%s/%s/_mlt?%s", index, _type, id, api.Pretty(pretty))
	body, err := api.DoCommand("GET", url, query)
	if err != nil {
		return retval, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}
	fmt.Println(body)
	return retval, err
}

type MoreLikeThisQuery struct {
	MoreLikeThis MLT `json:"more_like_this"`
}

type MLT struct {
	Fields              []string `json:"fields"`
	LikeText            string   `json:"like_text"`
	PercentTermsToMatch float32  `json:"percent_terms_to_match"`
	MinTermFrequency    int      `json:"min_term_freq"`
	MaxQueryTerms       int      `json:"max_query_terms"`
	StopWords           []string `json:"stop_words"`
	MinDocFrequency     int      `json:"min_doc_freq"`
	MaxDocFrequency     int      `json:"max_doc_freq"`
	MinWordLength       int      `json:"min_word_len"`
	MaxWordLength       int      `json:"max_word_len"`
	BoostTerms          int      `json:"boost_terms"`
	Boost               float32  `json:"boost"`
	Analyzer            string   `json:"analyzer"`
}
