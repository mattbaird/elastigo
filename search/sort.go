package search

import (
	"encoding/json"
	"fmt"
)

// Sorting accepts any number of Sort commands
//
//     Query().Sort(
//         Sort("last_name").Desc(),
//         Sort("age"),
//     )
func Sort(field string) *SortDsl {
	return &SortDsl{Name: field}
}

type SortBody []interface{}
type SortDsl struct {
	Name   string
	IsDesc bool
}

func (s *SortDsl) Desc() *SortDsl {
	s.IsDesc = true
	return s
}
func (s *SortDsl) Asc() *SortDsl {
	s.IsDesc = false
	return s
}

func (s *SortDsl) MarshalJSON() ([]byte, error) {
	if s.IsDesc {
		return json.Marshal(map[string]string{s.Name: "desc"})
	}
	if s.Name == "_score" {
		return []byte(`"_score"`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, s.Name)), nil // "user"  assuming default = asc?
}
