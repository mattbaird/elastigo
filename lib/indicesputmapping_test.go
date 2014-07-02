// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elastigo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func setup(t *testing.T) *Conn {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	c := NewConn()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	c.Domain = strings.Split(serverURL.Host, ":")[0]
	c.Port = strings.Split(serverURL.Host, ":")[1]

	return c
}

func teardown() {
	server.Close()
}

type TestStruct struct {
	Id            string `json:"id" elastic:"index:not_analyzed"`
	DontIndex     string `json:"dontIndex" elastic:"index:no"`
	Number        int    `json:"number" elastic:"type:integer,index:analyzed"`
	Omitted       string `json:"-"`
	NoJson        string `elastic:"type:string"`
	unexported    string
	JsonOmitEmpty string `json:"jsonOmitEmpty,omitempty" elastic:"type:string"`
	Embedded
	Nested       NestedStruct   `json:"nested"`
	NestedP      *NestedStruct  `json:"pointer_to_nested"`
	NestedS      []NestedStruct `json:"slice_of_nested"`
	MultiAnalyze string         `json:"multi_analyze"`
}

type Embedded struct {
	EmbeddedField string `json:"embeddedField" elastic:"type:string"`
}

type NestedStruct struct {
	NestedField string `json:"nestedField" elastic:"type:date"`
}

func TestPutMapping(t *testing.T) {
	c := setup(t)
	defer teardown()

	options := MappingOptions{
		Timestamp: TimestampOptions{Enabled: true},
		Id:        IdOptions{Index: "analyzed", Path: "id"},
		Properties: map[string]interface{}{
			// special properties that can't be expressed as tags
			"multi_analyze": map[string]interface{}{
				"type": "multi_field",
				"fields": map[string]map[string]string{
					"ma_analyzed":    {"type": "string", "index": "analyzed"},
					"ma_notanalyzed": {"type": "string", "index": "not_analyzed"},
				},
			},
		},
	}
	expValue := MappingForType("myType", MappingOptions{
		Timestamp: TimestampOptions{Enabled: true},
		Id:        IdOptions{Index: "analyzed", Path: "id"},
		Properties: map[string]interface{}{
			"NoJson":        map[string]string{"type": "string"},
			"dontIndex":     map[string]string{"index": "no"},
			"embeddedField": map[string]string{"type": "string"},
			"id":            map[string]string{"index": "not_analyzed"},
			"jsonOmitEmpty": map[string]string{"type": "string"},
			"number":        map[string]string{"index": "analyzed", "type": "integer"},
			"multi_analyze": map[string]interface{}{
				"type": "multi_field",
				"fields": map[string]map[string]string{
					"ma_analyzed":    {"type": "string", "index": "analyzed"},
					"ma_notanalyzed": {"type": "string", "index": "not_analyzed"},
				},
			},
			"nested": map[string]map[string]map[string]string{
				"properties": {
					"nestedField": {"type": "date"},
				},
			},
			"pointer_to_nested": map[string]map[string]map[string]string{
				"properties": {
					"nestedField": {"type": "date"},
				},
			},
			"slice_of_nested": map[string]map[string]map[string]string{
				"properties": {
					"nestedField": {"type": "date"},
				},
			},
		},
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var value map[string]interface{}
		bd, err := ioutil.ReadAll(r.Body)
		json.NewDecoder(strings.NewReader(string(bd))).Decode(&value)
		expValJson, err := json.MarshalIndent(expValue, "", "  ")
		if err != nil {
			t.Errorf("Got error: %v", err)
		}
		valJson, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			t.Errorf("Got error: %v", err)
		}

		if string(expValJson) != string(valJson) {
			t.Errorf("Expected %s but got %s", string(expValJson), string(valJson))
		}
	})

	err := c.PutMapping("myIndex", "myType", TestStruct{}, options)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
