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

package core

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// Index adds or updates a typed JSON document in a specific index, making it searchable, creating an index
// if it did not exist.
// if id is omited, op_type 'create' will be pased and http method will default to "POST"
// id is optional
// parentId is optional
// version is optional
// op_type is optional
// routing is optional
// timestamp is optional
// ttl is optional
// percolate is optional
// timeout is optional
// http://www.elasticsearch.org/guide/reference/api/index_.html
func Index(index string, _type string, id string, args map[string]interface{}, data interface{}) (api.BaseResponse, error) {
	var retval api.BaseResponse
	var url string

	if len(id) > 0 {
		url = fmt.Sprintf("/%s/%s/%s", index, _type, id)
	} else {
		url = fmt.Sprintf("/%s/%s", index, _type)
	}

	var method string
	if len(id) == 0 {
		method = "POST"
	} else {
		method = "PUT"
	}

	body, err := api.DoCommand(method, url, args, data)
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
	return retval, err
}
