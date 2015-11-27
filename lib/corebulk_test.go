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
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	"strings"

	"github.com/araddon/gou"
	"github.com/bmizerany/assert"
)

//  go test -bench=".*"
//  go test -bench="Bulk"

type sharedBuffer struct {
	mu     sync.Mutex
	Buffer []*bytes.Buffer
}

func NewSharedBuffer() *sharedBuffer {
	return &sharedBuffer{
		Buffer: make([]*bytes.Buffer, 0),
	}
}

func (b *sharedBuffer) Append(buf *bytes.Buffer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Buffer = append(b.Buffer, buf)
}

func (b *sharedBuffer) Length() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.Buffer)
}

func init() {
	flag.Parse()
	if testing.Verbose() {
		gou.SetupLogging("debug")
	}
}

// take two ints, compare, need to be within 5%
func closeInt(a, b int) bool {
	c := float64(a) / float64(b)
	if c >= .95 && c <= 1.05 {
		return true
	}
	return false
}

func TestBulkIndexerBasic(t *testing.T) {
	testIndex := "users"
	var (
		buffers        = NewSharedBuffer()
		totalBytesSent int
		messageSets    int
	)

	InitTests(true)
	c := NewTestConn()

	c.DeleteIndex(testIndex)

	indexer := c.NewBulkIndexer(3)
	indexer.Sender = func(buf *bytes.Buffer) error {
		messageSets += 1
		totalBytesSent += buf.Len()
		buffers.Append(buf)
		//log.Printf("buffer:%s", string(buf.Bytes()))
		return indexer.Send(buf)
	}
	indexer.Start()

	date := time.Unix(1257894000, 0)
	data := map[string]interface{}{
		"name": "smurfs",
		"age":  22,
		"date": "yesterday",
	}

	err := indexer.Index(testIndex, "user", "1", "", "", &date, nil, data)
	waitFor(func() bool {
		return buffers.Length() > 0
	}, 5)

	// part of request is url, so lets factor that in
	//totalBytesSent = totalBytesSent - len(*eshost)
	assert.T(t, buffers.Length() == 1, fmt.Sprintf("Should have sent one operation but was %d", buffers.Length()))
	assert.T(t, indexer.NumErrors() == 0 && err == nil, fmt.Sprintf("Should not have any errors. NumErrors: %v, err: %v", indexer.NumErrors(), err))
	expectedBytes := 129
	assert.T(t, totalBytesSent == expectedBytes, fmt.Sprintf("Should have sent %v bytes but was %v", expectedBytes, totalBytesSent))

	err = indexer.Index(testIndex, "user", "2", "", "", nil, nil, data)
	waitFor(func() bool {
		return buffers.Length() > 1
	}, 5)

	// this will test to ensure that Flush actually catches a doc
	indexer.Flush()
	<-time.After(time.Millisecond * 1) // we need to wait for the httpSender to read from the send buffer
	totalBytesSent = totalBytesSent - len(*eshost)
	assert.T(t, err == nil, fmt.Sprintf("Should have nil error  =%v", err))
	assert.T(t, buffers.Length() == 2, fmt.Sprintf("Should have another buffer ct=%d", buffers.Length()))

	assert.T(t, indexer.NumErrors() == 0, fmt.Sprintf("Should not have any errors %d", indexer.NumErrors()))
	expectedBytes = 220
	assert.T(t, closeInt(totalBytesSent, expectedBytes), fmt.Sprintf("Should have sent %v bytes but was %v", expectedBytes, totalBytesSent))

	indexer.Stop()
}

func TestBulkIndexerErrors(t *testing.T) {
	testIndex := "users"
	InitTests(true)
	c := NewTestConn()

	c.DeleteIndex(testIndex)

	sent := make(chan struct{}, 1)

	indexer := c.NewBulkIndexer(3)
	indexer.Sender = func(buf *bytes.Buffer) error {
		// log.Printf("buffer:%s", string(buf.Bytes()))
		ret := indexer.Send(buf)
		sent <- struct{}{}
		return ret
	}
	indexer.Start()
	errch := make(chan *ErrorBuffer, 1)
	indexer.ErrorChannel = errch

	date := time.Unix(1257894000, 0)

	data := map[string]interface{}{
		"name": "smurfettes",
		"age":  21,
		"date": "today",
	}
	data2 := map[string]interface{}{
		"name": "smurfs",
		"age":  "this is not an int",
		"date": "yesterday",
	}

	_, err := c.DoCommand("PUT", fmt.Sprintf("/%s", testIndex), nil,
		`{
	  "mappings": {
		"user": {
		  "properties": {
			"age": { "type": "integer" }
		  }
		}
	  }
	}`)
	assert.T(t, err == nil, fmt.Sprintf("Should not return an error: %v", err))

	//act
	err = indexer.Index(testIndex, "user", "1", "", "", &date, nil, data)
	err2 := indexer.Index(testIndex, "user", "2", "", "", &date, nil, data2)

	<-sent
	//assert
	assert.T(t, err == nil, fmt.Sprintf("Should not return an error: %v", err))
	assert.T(t, err2 == nil, fmt.Sprintf("Should not return an error: %v", err2))
	time.Sleep(1 * time.Microsecond)
	assert.T(t, indexer.NumErrors() == 1, fmt.Sprintf("Should have recorded 1 error but saw: %d", indexer.NumErrors()))

	errBuf := <-errch
	bulkErr, ok := errBuf.Err.(BulkIndexingError)
	assert.T(t, ok, fmt.Sprintf("Error should have been a BulkIndexingError but was %T", errBuf.Err))

	assert.T(t, len(bulkErr.Items) == 2, fmt.Sprintf("Expected 2 items in response but got: %d", len(bulkErr.Items)))
	status1 := int(bulkErr.Items[1]["index"].(map[string]interface{})["status"].(float64))
	assert.T(t, status1 == 400, fmt.Sprintf("Expected second item to have status 400 but was %d", status1))
	status0 := int(bulkErr.Items[0]["index"].(map[string]interface{})["status"].(float64))
	assert.T(t, status0 == 201, fmt.Sprintf("Expected first item to have status 201 but was %d", status0))

	lines := strings.Split(errBuf.Buf.String(), "\n")
	assert.T(t, lines[0] == `{"index":{"_index":"users","_type":"user","_id":"1","_timestamp":"1257894000000"}}`, fmt.Sprintf("Expected index header but got: %s", lines[0]))
	assert.T(t, lines[1] == `{"age":21,"date":"today","name":"smurfettes"}`, fmt.Sprintf("Expected document but got: %s", lines[1]))
	assert.T(t, lines[2] == `{"index":{"_index":"users","_type":"user","_id":"2","_timestamp":"1257894000000"}}`, fmt.Sprintf("Expected index header but got: %s", lines[2]))
	assert.T(t, lines[3] == `{"age":"this is not an int","date":"yesterday","name":"smurfs"}`, fmt.Sprintf("Expected document but got: %s", lines[3]))
}

func TestRefreshParam(t *testing.T) {
	requrlChan := make(chan *url.URL, 1)
	InitTests(true)
	c := NewTestConn()
	c.RequestTracer = func(method, urlStr, body string) {
		requrl, _ := url.Parse(urlStr)
		requrlChan <- requrl
	}
	date := time.Unix(1257894000, 0)
	data := map[string]interface{}{"name": "smurfs", "age": 22, "date": date}

	// Now tests small batches
	indexer := c.NewBulkIndexer(1)
	indexer.Refresh = true

	indexer.Start()
	<-time.After(time.Millisecond * 20)

	indexer.Index("users", "user", "2", "", "", &date, nil, data)

	<-time.After(time.Millisecond * 200)
	//	indexer.Flush()
	indexer.Stop()
	requrl := <-requrlChan
	assert.T(t, requrl.Query().Get("refresh") == "true", "Should have set refresh query param to true")
}

func TestWithoutRefreshParam(t *testing.T) {
	requrlChan := make(chan *url.URL, 1)
	InitTests(true)
	c := NewTestConn()
	c.RequestTracer = func(method, urlStr, body string) {
		requrl, _ := url.Parse(urlStr)
		requrlChan <- requrl
	}
	date := time.Unix(1257894000, 0)
	data := map[string]interface{}{"name": "smurfs", "age": 22, "date": date}

	// Now tests small batches
	indexer := c.NewBulkIndexer(1)

	indexer.Start()
	<-time.After(time.Millisecond * 20)

	indexer.Index("users", "user", "2", "", "", &date, nil, data)

	<-time.After(time.Millisecond * 200)
	//	indexer.Flush()
	indexer.Stop()
	requrl := <-requrlChan
	assert.T(t, requrl.Query().Get("refresh") == "false", "Should have set refresh query param to false")
}

// currently broken in drone.io
func XXXTestBulkUpdate(t *testing.T) {
	var (
		buffers        = NewSharedBuffer()
		totalBytesSent int
		messageSets    int
	)

	InitTests(true)
	c := NewTestConn()
	c.Port = "9200"
	indexer := c.NewBulkIndexer(3)
	indexer.Sender = func(buf *bytes.Buffer) error {
		messageSets += 1
		totalBytesSent += buf.Len()
		buffers.Append(buf)
		return indexer.Send(buf)
	}
	indexer.Start()

	date := time.Unix(1257894000, 0)
	user := map[string]interface{}{
		"name": "smurfs", "age": 22, "date": date, "count": 1,
	}

	// Lets make sure the data is in the index ...
	_, err := c.Index("users", "user", "5", nil, user)

	// script and params
	data := map[string]interface{}{
		"script": "ctx._source.count += 2",
	}
	err = indexer.Update("users", "user", "5", "", "", &date, nil, data)
	// So here's the deal. Flushing does seem to work, you just have to give the
	// channel a moment to recieve the message ...
	//	<- time.After(time.Millisecond * 20)
	//	indexer.Flush()

	waitFor(func() bool {
		return buffers.Length() > 0
	}, 5)

	indexer.Stop()

	assert.T(t, indexer.NumErrors() == 0 && err == nil, fmt.Sprintf("Should not have any errors, bulkErrorCt:%v, err:%v", indexer.NumErrors(), err))

	response, err := c.Get("users", "user", "5", nil)
	assert.T(t, err == nil, fmt.Sprintf("Should not have any errors  %v", err))
	m := make(map[string]interface{})
	json.Unmarshal([]byte(*response.Source), &m)
	newCount := m["count"]
	assert.T(t, newCount.(float64) == 3,
		fmt.Sprintf("Should have update count: %#v ... %#v", m["count"], response))
}

func TestBulkSmallBatch(t *testing.T) {
	var (
		messageSets int
	)

	InitTests(true)
	c := NewTestConn()

	date := time.Unix(1257894000, 0)
	data := map[string]interface{}{"name": "smurfs", "age": 22, "date": date}

	// Now tests small batches
	indexer := c.NewBulkIndexer(1)
	indexer.BufferDelayMax = 100 * time.Millisecond
	indexer.BulkMaxDocs = 2
	messageSets = 0
	indexer.Sender = func(buf *bytes.Buffer) error {
		messageSets += 1
		return indexer.Send(buf)
	}
	indexer.Start()
	<-time.After(time.Millisecond * 20)

	indexer.Index("users", "user", "2", "", "", &date, nil, data)
	indexer.Index("users", "user", "3", "", "", &date, nil, data)
	indexer.Index("users", "user", "4", "", "", &date, nil, data)
	<-time.After(time.Millisecond * 200)
	//	indexer.Flush()
	indexer.Stop()
	assert.T(t, messageSets == 2, fmt.Sprintf("Should have sent 2 message sets %d", messageSets))

}

func TestBulkDelete(t *testing.T) {
	InitTests(true)
	var lock sync.Mutex
	c := NewTestConn()
	indexer := c.NewBulkIndexer(1)
	sentBytes := []byte{}

	indexer.Sender = func(buf *bytes.Buffer) error {
		lock.Lock()
		sentBytes = append(sentBytes, buf.Bytes()...)
		lock.Unlock()
		return nil
	}

	indexer.Start()

	indexer.Delete("fake", "fake_type", "1", nil)

	indexer.Flush()
	indexer.Stop()

	lock.Lock()
	sent := string(sentBytes)
	lock.Unlock()

	expected := `{"delete":{"_index":"fake","_type":"fake_type","_id":"1"}}
`
	asExpected := sent == expected
	assert.T(t, asExpected, fmt.Sprintf("Should have sent '%s' but actually sent '%s'", expected, sent))
}

func XXXTestBulkErrors(t *testing.T) {
	// lets set a bad port, and hope we get a conn refused error?
	c := NewTestConn()
	c.Port = "27845"
	defer func() {
		c.Port = "9200"
	}()
	indexer := c.NewBulkIndexerErrors(10, 1)
	indexer.Start()
	errorCt := 0
	go func() {
		for i := 0; i < 20; i++ {
			date := time.Unix(1257894000, 0)
			data := map[string]interface{}{"name": "smurfs", "age": 22, "date": date}
			indexer.Index("users", "user", strconv.Itoa(i), "", "", &date, nil, data)
		}
	}()
	var errBuf *ErrorBuffer
	for errBuf = range indexer.ErrorChannel {
		errorCt++
		break
	}
	if errBuf.Buf.Len() > 0 {
		gou.Debug(errBuf.Err)
	}
	assert.T(t, errorCt > 0, fmt.Sprintf("ErrorCt should be > 0 %d", errorCt))
	indexer.Stop()
}

func TestBulkVersioning_Internal(t *testing.T) {
	testIndex := "users"
	var (
		buffers        = NewSharedBuffer()
		totalBytesSent int
		messageSets    int
	)

	InitTests(true)
	c := NewTestConn()
	c.RequestTracer = func(method, url, body string) {
		t.Logf("%s %s HTTP/1.1\n%s", method, url, body)
	}

	c.DeleteIndex(testIndex)

	indexer := c.NewBulkIndexer(3)
	indexer.Sender = func(buf *bytes.Buffer) error {
		messageSets += 1
		totalBytesSent += buf.Len()
		buffers.Append(buf)
		// log.Printf("buffer:%s", string(buf.Bytes()))
		return indexer.Send(buf)
	}
	errCh := make(chan *ErrorBuffer)
	indexer.ErrorChannel = errCh
	indexer.Start()

	date := time.Unix(1257894000, 0)
	data := map[string]interface{}{
		"name": "smurfs",
		"age":  22,
		"date": "yesterday",
	}

	indexer.Index(testIndex, "user", "1", "", "", &date, nil, data)

	//act
	data["extra"] = "1"
	indexer.Index(testIndex, "user", "1", "", "", &date, &DocVersion{Version: 1}, data)

	indexer.Update(testIndex, "user", "1", "", "", &date, &DocVersion{Version: 2, VersionType: "internal"}, map[string]interface{}{
		"doc": map[string]interface{}{
			"updated": true,
		},
	})

	data["extra"] = "3"
	indexer.Delete(testIndex, "user", "1", &DocVersion{Version: 7})

	//assert
	errBuf := <-errCh

	bulkErr, ok := errBuf.Err.(BulkIndexingError)
	assert.T(t, ok, fmt.Sprintf("Expected bulk indexing error but was: %T\n\t%v", errBuf.Err, errBuf.Err))

	js, _ := json.MarshalIndent(bulkErr.Items, "", "  ")
	t.Logf("Items:%s", string(js))

	assert.T(t, getStatus(0, bulkErr.Items) == 201, "First should be created")
	assert.T(t, getVersion(0, bulkErr.Items) == int64(1), "First should have version 1")
	assert.T(t, getStatus(1, bulkErr.Items) == 200, "Should be reindexed with version 2")
	assert.T(t, getVersion(1, bulkErr.Items) == int64(2), "Should be reindexed with version 2")
	assert.T(t, getStatus(2, bulkErr.Items) == 200, "Should be updated with version 3")
	assert.T(t, getVersion(2, bulkErr.Items) == int64(3), "Should be updated with version 3")
	assert.T(t, getStatus(3, bulkErr.Items) == 409, "Should fail to delete due to version conflict")
	assert.T(t, getError(3, bulkErr.Items) == "VersionConflictEngineException[[users][2] [user][1]: version conflict, current [3], provided [7]]")

}

func TestBulkVersioning_External(t *testing.T) {
	testIndex := "users"
	var (
		buffers        = NewSharedBuffer()
		totalBytesSent int
		messageSets    int
	)

	InitTests(true)
	c := NewTestConn()
	c.RequestTracer = func(method, url, body string) {
		t.Logf("%s %s HTTP/1.1\n%s", method, url, body)
	}

	c.DeleteIndex(testIndex)

	indexer := c.NewBulkIndexer(3)
	indexer.Sender = func(buf *bytes.Buffer) error {
		messageSets += 1
		totalBytesSent += buf.Len()
		buffers.Append(buf)
		// log.Printf("buffer:%s", string(buf.Bytes()))
		return indexer.Send(buf)
	}
	errCh := make(chan *ErrorBuffer)
	indexer.ErrorChannel = errCh
	indexer.Start()

	date := time.Unix(1257894000, 0)
	data := map[string]interface{}{
		"name": "smurfs",
		"age":  22,
		"date": "yesterday",
	}

	now := time.Now().Unix()
	indexer.Index(testIndex, "user", "1", "", "", &date, EXTERNAL.V(now), data)

	//act
	data["extra"] = "1"
	indexer.Index(testIndex, "user", "1", "", "", &date, EXTERNAL_GT.V(now), data)

	data["extra"] = "2"
	indexer.Index(testIndex, "user", "1", "", "", &date, EXTERNAL_GT.V(now+2), data)

	data["extra"] = "3"
	indexer.Delete(testIndex, "user", "1", EXTERNAL_GT.V(now-1))

	//assert
	errBuf := <-errCh

	bulkErr, ok := errBuf.Err.(BulkIndexingError)
	assert.T(t, ok, fmt.Sprintf("Expected bulk indexing error but was: %T\n\t%v", errBuf.Err, errBuf.Err))

	js, _ := json.MarshalIndent(bulkErr.Items, "", "  ")
	t.Logf("Items:%s", string(js))

	assert.T(t, getStatus(0, bulkErr.Items) == 201, "First should be created")
	assert.T(t, getVersion(0, bulkErr.Items) == now, "First should have version now")
	assert.T(t, getStatus(1, bulkErr.Items) == 409, "Should fail to reindex with same version")
	assert.T(t, getStatus(2, bulkErr.Items) == 200, "Should be updated with version now+2")
	assert.T(t, getVersion(2, bulkErr.Items) == now+2, "Should be updated with version now+2")
	assert.T(t, getStatus(3, bulkErr.Items) == 409, "Should fail to delete due to version conflict")

}

func getStatus(index int, items []map[string]interface{}) int {
	for _, doc := range items[index] {
		doc := doc.(map[string]interface{})
		return int(doc["status"].(float64))
	}
	panic(fmt.Sprintf("no properties in %v", items[index]))
}

func getVersion(index int, items []map[string]interface{}) int64 {
	for _, doc := range items[index] {
		doc := doc.(map[string]interface{})
		return int64(doc["_version"].(float64))
	}
	panic(fmt.Sprintf("no properties in %v", items[index]))
}

func getError(index int, items []map[string]interface{}) string {
	for _, doc := range items[index] {
		doc := doc.(map[string]interface{})
		return doc["error"].(string)
	}
	panic(fmt.Sprintf("no properties in %v", items[index]))
}

/*
BenchmarkSend	18:33:00 bulk_test.go:131: Sent 1 messages in 0 sets totaling 0 bytes
18:33:00 bulk_test.go:131: Sent 100 messages in 1 sets totaling 145889 bytes
18:33:01 bulk_test.go:131: Sent 10000 messages in 100 sets totaling 14608888 bytes
18:33:05 bulk_test.go:131: Sent 20000 messages in 99 sets totaling 14462790 bytes
   20000	    234526 ns/op

*/
func BenchmarkSend(b *testing.B) {
	InitTests(true)
	c := NewTestConn()
	b.StartTimer()
	totalBytes := 0
	sets := 0
	indexer := c.NewBulkIndexer(1)
	indexer.Sender = func(buf *bytes.Buffer) error {
		totalBytes += buf.Len()
		sets += 1
		//log.Println("got bulk")
		return indexer.Send(buf)
	}
	for i := 0; i < b.N; i++ {
		about := make([]byte, 1000)
		rand.Read(about)
		data := map[string]interface{}{"name": "smurfs", "age": 22, "date": time.Unix(1257894000, 0), "about": about}
		indexer.Index("users", "user", strconv.Itoa(i), "", "", nil, nil, data)
	}
	log.Printf("Sent %d messages in %d sets totaling %d bytes \n", b.N, sets, totalBytes)
	if indexer.NumErrors() != 0 {
		b.Fail()
	}
}

/*
TODO:  this should be faster than above

BenchmarkSendBytes	18:33:05 bulk_test.go:169: Sent 1 messages in 0 sets totaling 0 bytes
18:33:05 bulk_test.go:169: Sent 100 messages in 2 sets totaling 292299 bytes
18:33:09 bulk_test.go:169: Sent 10000 messages in 99 sets totaling 14473800 bytes
   10000	    373529 ns/op

*/
func BenchmarkSendBytes(b *testing.B) {
	InitTests(true)
	c := NewTestConn()
	about := make([]byte, 1000)
	rand.Read(about)
	data := map[string]interface{}{"name": "smurfs", "age": 22, "date": time.Unix(1257894000, 0), "about": about}
	body, _ := json.Marshal(data)
	b.StartTimer()
	totalBytes := 0
	sets := 0
	indexer := c.NewBulkIndexer(1)
	indexer.Sender = func(buf *bytes.Buffer) error {
		totalBytes += buf.Len()
		sets += 1
		return indexer.Send(buf)
	}
	for i := 0; i < b.N; i++ {
		indexer.Index("users", "user", strconv.Itoa(i), "", "", nil, nil, body)
	}
	log.Printf("Sent %d messages in %d sets totaling %d bytes \n", b.N, sets, totalBytes)
	if indexer.NumErrors() != 0 {
		b.Fail()
	}
}
