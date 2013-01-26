/**
   Copyright 2012 Matthew Baird

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
**/
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
)

type Request http.Request

const (
	Version         = "0.0.1"
	DefaultProtocol = "http"
	DefaultDomain   = "localhost"
	DefaultPort     = "9200"
)

var (
	_               = log.Ldate
	Protocol string = DefaultProtocol
	Domain   string = DefaultDomain
	Port     string = DefaultPort
)

func ElasticSearchRequest(method, path string) (*Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s://%s:%s%s", Protocol, Domain, Port, path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "elasticSearch/"+Version+" ("+runtime.GOOS+"-"+runtime.GOARCH+")")
	return (*Request)(req), nil
}

func (r *Request) SetBodyJson(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.SetBody(bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	return nil
}

func (r *Request) SetBodyString(body string) {
	r.SetBody(strings.NewReader(body))
}

func (r *Request) SetBody(body io.Reader) {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *strings.Reader:
			r.ContentLength = int64(v.Len())
		case *bytes.Buffer:
			r.ContentLength = int64(v.Len())
		}
	}
}

func (r *Request) Do(v interface{}) (int, []byte, error) {
	res, err := http.DefaultClient.Do((*http.Request)(r))
	if err != nil {
		return 0, nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)

	if res.StatusCode > 304 && v != nil {
		jsonErr := json.Unmarshal(bodyBytes, v)
		if jsonErr != nil {
			return res.StatusCode, nil, err
		}
	}
	return res.StatusCode, bodyBytes, nil
}
