package core

import (
	"log"
	"testing"
)

/*
TODO
  - load test data
  - 
*/
func TestDsl(t *testing.T) {
	out, err := Search("logstash-2012.12.29").Pretty().Facet(
		Facet().Fields("@fields.category").Size("25"),
	).Query(
		Query().All(),
	).Result()
	//log.Println(out)
	if out != nil {
		log.Println(string(out.Facets))
	}
	log.Println(err)

	out, err = Search("logstash-2012.12.29").Pretty().Facet(
		Facet().Fields("@fields.category").Size("25"),
	).Query(
		Query().Range(
			Range().Field("@timestamp").From("2012-12-29T16:52:48+00:00").To("2012-12-29T17:52:48+00:00"),
		).Search("player"),
	).Result()
	//log.Println(out)
	if out != nil {
		log.Println(string(out.Facets))
	}
	log.Println(err)
}
