package search

import (
	"encoding/json"
	"time"
)

func Aggregate(name string) *AggregateDsl {
	return &AggregateDsl{Name: name}
}

type AggregateDsl struct {
	Name          string
	TypeName      string
	Type          interface{}
	Filters       *FilterWrap              `json:"filters,omitempty"`
	AggregatesVal map[string]*AggregateDsl `json:"aggregations,omitempty"`
}

type FieldAggregate struct {
	Field string `json:"field"`
}

/**
 * Aggregates accepts n "sub-aggregates" to be applied to this aggregate
 *
 * agg := Aggregate("user").Term("user_id")
 * agg.Aggregates(
 *   Aggregate("total_spent").Sum("price"),
 *   Aggregate("total_saved").Sum("discount"),
 * )
 */
func (d *AggregateDsl) Aggregates(aggs ...*AggregateDsl) *AggregateDsl {
	if len(aggs) < 1 {
		return d
	}
	if len(d.AggregatesVal) == 0 {
		d.AggregatesVal = make(map[string]*AggregateDsl)
	}

	for _, agg := range aggs {
		d.AggregatesVal[agg.Name] = agg
	}
	return d
}

func (d *AggregateDsl) Min(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "min"
	return d
}

func (d *AggregateDsl) Max(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "max"
	return d
}

func (d *AggregateDsl) Sum(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "sum"
	return d
}

func (d *AggregateDsl) Avg(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "avg"
	return d
}

func (d *AggregateDsl) Stats(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "stats"
	return d
}

func (d *AggregateDsl) ExtendedStats(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "extended_stats"
	return d
}

func (d *AggregateDsl) ValueCount(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "value_count"
	return d
}

func (d *AggregateDsl) Percentiles(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "percentiles"
	return d
}

type Cardinality struct {
	Field              string  `json:"field"`
	PrecisionThreshold float64 `json:"precision_threshold,omitempty"`
	Rehash             bool    `json:"rehash,omitempty"`
}

/**
 * Cardinality(
 *	 "field_name",
 *	 true,
 *   0,
 * )
 */
func (d *AggregateDsl) Cardinality(field string, rehash bool, threshold int) *AggregateDsl {
	c := Cardinality{Field: field}

	// Only set if it's false, since the default is true
	if !rehash {
		c.Rehash = false
	}

	if threshold > 0 {
		c.PrecisionThreshold = float64(threshold)
	}
	d.Type = c
	d.TypeName = "cardinality"
	return d
}

func (d *AggregateDsl) Global() *AggregateDsl {
	d.Type = struct{}{}
	d.TypeName = "global"
	return d
}

func (d *AggregateDsl) Filter(filters ...interface{}) *AggregateDsl {

	if len(filters) == 0 {
		return d
	}

	if d.Filters == nil {
		d.Filters = NewFilterWrap()
	}

	d.Filters.addFilters(filters)
	return d
}

func (d *AggregateDsl) Missing(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "missing"
	return d
}

func (d *AggregateDsl) Terms(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "terms"
	return d
}

func (d *AggregateDsl) SignificantTerms(field string) *AggregateDsl {
	d.Type = FieldAggregate{Field: field}
	d.TypeName = "significant_terms"
	return d
}

type Histogram struct {
	Field          string      `json:"field"`
	Interval       float64     `json:"interval"`
	MinDocCount    float64     `json:"min_doc_count"`
	ExtendedBounds interface{} `json:"extended_bounds,omitempty"`
}

func (d *AggregateDsl) Histogram(field string, interval int) *AggregateDsl {
	d.Type = Histogram{
		Field:       field,
		Interval:    float64(interval),
		MinDocCount: 1,
	}
	d.TypeName = "histogram"
	return d
}

type DateHistogram struct {
	Field          string      `json:"field"`
	Interval       string      `json:"interval"`
	MinDocCount    float64     `json:"min_doc_count"`
	ExtendedBounds interface{} `json:"extended_bounds,omitempty"`
}

func (d *AggregateDsl) DateHistogram(field, interval string) *AggregateDsl {
	d.Type = DateHistogram{
		Field:       field,
		Interval:    interval,
		MinDocCount: 1,
	}
	d.TypeName = "date_histogram"
	return d
}

// Sets the min doc count for a date histogram or histogram
// This will no-op if used on an inappropriate dsl type
func (d *AggregateDsl) MinDocCount(i float64) *AggregateDsl {

	if d.TypeName == "date_histogram" {
		t := d.Type.(DateHistogram)
		t.MinDocCount = i
		d.Type = t
	} else if d.TypeName == "histogram" {
		t := d.Type.(Histogram)
		t.MinDocCount = i
		d.Type = t
	}

	return d
}

// Hackety hack function that expects different types depending on the type of aggregate
// Not very idiomatic, but fits the elastigo DSL
func (d *AggregateDsl) ExtendedBounds(min, max interface{}) *AggregateDsl {
	if min == nil && max == nil {
		return d
	}

	if d.TypeName == "date_histogram" {
		var n time.Time
		var x time.Time
		t := d.Type.(DateHistogram)
		if min != nil {
			switch min.(type) {
			case time.Time:
				n = min.(time.Time)
			}
		}
		if max != nil {
			switch max.(type) {
			case time.Time:
				x = max.(time.Time)
			}
		}

		if min == nil {
			bounds := struct {
				Max time.Time `json:"max"`
			}{x}
			t.ExtendedBounds = &bounds
		} else if max == nil {
			bounds := struct {
				Min time.Time `json:"min"`
			}{n}
			t.ExtendedBounds = &bounds
		} else {
			bounds := struct {
				Min time.Time `json:"min"`
				Max time.Time `json:"max"`
			}{n, x}
			t.ExtendedBounds = &bounds
		}

		d.Type = t
	}
	if d.TypeName == "histogram" {
		var n float64
		var x float64
		t := d.Type.(Histogram)
		if min != nil {
			switch min.(type) {
			case time.Time:
				n = min.(float64)
			}
		}
		if max != nil {
			switch max.(type) {
			case time.Time:
				x = max.(float64)
			}
		}

		if min == nil {
			bounds := struct {
				Max float64 `json:"max"`
			}{x}
			t.ExtendedBounds = &bounds
		} else if max == nil {
			bounds := struct {
				Min float64 `json:"min"`
			}{n}
			t.ExtendedBounds = &bounds
		} else {
			bounds := struct {
				Min float64 `json:"min"`
				Max float64 `json:"max"`
			}{n, x}
			t.ExtendedBounds = &bounds
		}

		d.Type = t
	}

	return d
}
func (d *AggregateDsl) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.toMap())
}

func (d *AggregateDsl) toMap() map[string]interface{} {
	root := map[string]interface{}{}

	if d.Type != nil {
		root[d.TypeName] = d.Type
	}
	aggregates := d.aggregatesMap()

	if d.Filters != nil {
		root["filter"] = d.Filters
	}

	if len(aggregates) > 0 {
		root["aggregations"] = aggregates
	}
	return root

}
func (d *AggregateDsl) aggregatesMap() map[string]interface{} {
	root := map[string]interface{}{}

	if len(d.AggregatesVal) > 0 {
		for _, agg := range d.AggregatesVal {
			root[agg.Name] = agg.toMap()
		}
	}
	return root
}
