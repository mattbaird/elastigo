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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	hostpool "github.com/bitly/go-hostpool"
)

type Request struct {
	*http.Client
	*http.Request
	hostResponse hostpool.HostPoolResponse
}

func (r *Request) SetBodyJson(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.SetBodyBytes(body)
	r.Header.Set("Content-Type", "application/json")
	return nil
}

func (r *Request) SetBodyString(body string) {
	r.SetBody(strings.NewReader(body))
}

func (r *Request) SetBodyBytes(body []byte) {
	r.SetBody(bytes.NewReader(body))
}

func (r *Request) SetBody(body io.Reader) {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.Body = rc
	r.ContentLength = -1
}

func (r *Request) Do(v interface{}) (int, []byte, error) {
	response, bodyBytes, err := r.DoResponse(v)
	if response == nil {
		return -1, bodyBytes, err
	}
	return response.StatusCode, bodyBytes, err
}

func (r *Request) DoResponse(v interface{}) (*http.Response, []byte, error) {
	var client = r.Client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(r.Request)
	// Inform the HostPool of what happened to the request and allow it to update
	r.hostResponse.Mark(err)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return res, nil, err
	}

	if res.StatusCode == 404 {
		return res, bodyBytes, RecordNotFound
	}

	if res.StatusCode > 304 && v != nil {
		err = json.Unmarshal(bodyBytes, v)
	}
	return res, bodyBytes, err
}

func Escape(args map[string]interface{}) (s string, err error) {
	vals := url.Values{}
	for key, val := range args {
		switch v := val.(type) {
		case string:
			vals.Add(key, v)
		case bool:
			vals.Add(key, strconv.FormatBool(v))
		case int, int32, int64:
			vInt := reflect.ValueOf(v).Int()
			vals.Add(key, strconv.FormatInt(vInt, 10))
		case float32, float64:
			vFloat := reflect.ValueOf(v).Float()
			vals.Add(key, strconv.FormatFloat(vFloat, 'f', -1, 32))
		case []string:
			vals.Add(key, strings.Join(v, ","))
		default:
			err = fmt.Errorf("Could not format URL argument: %s", key)
			return
		}
	}
	s = vals.Encode()
	return
}
