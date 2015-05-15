package elastigo

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidIndexLine = errors.New("Cannot parse indexline")

// Create an IndexInfo from the string _cat/indices would produce
func NewCatIndexInfo(indexLine string) (catIndex *CatIndexInfo, err error) {
	split := strings.Fields(indexLine)
	if len(split) < 4 {
		return nil, ErrInvalidIndexLine
	}
	catIndex = &CatIndexInfo{}
	catIndex.Store = CatIndexStore{}
	catIndex.Docs = CatIndexDocs{}
	catIndex.Health = split[0]
	catIndex.Name = split[1]
	catIndex.Shards, err = strconv.Atoi(split[2])
	if err != nil {
		catIndex.Shards = 0
	}
	catIndex.Replicas, err = strconv.Atoi(split[3])
	if err != nil {
		catIndex.Replicas = 0
	}
	if len(split) == 4 {
		return catIndex, nil
	}
	catIndex.Docs.Count, err = strconv.ParseInt(split[4], 10, 64)
	if err != nil {
		catIndex.Docs.Count = 0
	}
	if len(split) == 5 {
		return catIndex, nil
	}
	catIndex.Docs.Deleted, err = strconv.ParseInt(split[5], 10, 64)
	if err != nil {
		catIndex.Docs.Deleted = 0
	}
	if len(split) == 6 {
		return catIndex, nil
	}
	catIndex.Store.Size, err = strconv.ParseInt(split[6], 10, 64)
	if err != nil {
		catIndex.Store.Size = 0
	}
	if len(split) == 7 {
		return catIndex, nil
	}
	catIndex.Store.PriSize, err = strconv.ParseInt(split[7], 10, 64)
	if err != nil {
		catIndex.Store.PriSize = 0
	}
	return catIndex, nil
}
