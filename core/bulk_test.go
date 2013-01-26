package core

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"log"
	"strconv"
	"testing"
	"time"
)

//  go test -bench=".*" 
//  go test -bench="Bulk" 

var (
	buffers        = make([]*bytes.Buffer, 0)
	totalBytesSent int
	messageSets    int
)

func TestBulk(t *testing.T) {
	InitTests(true)
	BulkSendor = func(buf *bytes.Buffer) {
		messageSets += 1
		totalBytesSent += buf.Len()
		buffers = append(buffers, buf)
		BulkSend(buf)
	}

	date := time.Unix(1257894000, 0)
	data := map[string]interface{}{"name": "smurfs", "age": 22, "date": time.Unix(1257894000, 0)}
	err := IndexBulk("users", "user", "1", &date, data)

	WaitFor(func() bool {
		return len(buffers) > 0
	}, 5)
	// part of request is url, so lets factor that in
	totalBytesSent = totalBytesSent - len(*eshost)
	Assert(len(buffers) == 1, t, "Should have sent one operation")
	Assert(BulkErrorCt == 0 && err == nil, t, "Should not have any errors")
	Assert(totalBytesSent == 135, t, "Should have sent 135 bytes but was %v", totalBytesSent)

	err = IndexBulk("users", "user", "2", nil, data)

	WaitFor(func() bool {
		return len(buffers) > 1
	}, 5)
	totalBytesSent = totalBytesSent - len(*eshost)
	Assert(len(buffers) == 2, t, "Should have nil error, and another buffer")

	Assert(BulkErrorCt == 0 && err == nil, t, "Should not have any errors")
	Assert(totalBytesSent == 241, t, "Should have sent 241 bytes but was %v", totalBytesSent)
}

/*
BenchmarkBulkSend	18:33:00 bulk_test.go:131: Sent 1 messages in 0 sets totaling 0 bytes 
18:33:00 bulk_test.go:131: Sent 100 messages in 1 sets totaling 145889 bytes 
18:33:01 bulk_test.go:131: Sent 10000 messages in 100 sets totaling 14608888 bytes 
18:33:05 bulk_test.go:131: Sent 20000 messages in 99 sets totaling 14462790 bytes 
   20000	    234526 ns/op

*/
func BenchmarkBulkSend(b *testing.B) {
	InitTests(true)
	b.StartTimer()
	totalBytes := 0
	sets := 0
	BulkSendor = func(buf *bytes.Buffer) {
		totalBytes += buf.Len()
		sets += 1
		//log.Println("got bulk")
		BulkSend(buf)
	}
	for i := 0; i < b.N; i++ {
		about := make([]byte, 1000)
		rand.Read(about)
		data := map[string]interface{}{"name": "smurfs", "age": 22, "date": time.Unix(1257894000, 0), "about": about}
		IndexBulk("users", "user", strconv.Itoa(i), nil, data)
	}
	log.Printf("Sent %d messages in %d sets totaling %d bytes \n", b.N, sets, totalBytes)
	if BulkErrorCt != 0 {
		b.Fail()
	}
}

/*
TODO:  this should be faster than above

BenchmarkBulkSendBytes	18:33:05 bulk_test.go:169: Sent 1 messages in 0 sets totaling 0 bytes 
18:33:05 bulk_test.go:169: Sent 100 messages in 2 sets totaling 292299 bytes 
18:33:09 bulk_test.go:169: Sent 10000 messages in 99 sets totaling 14473800 bytes 
   10000	    373529 ns/op

*/
func BenchmarkBulkSendBytes(b *testing.B) {
	InitTests(true)
	about := make([]byte, 1000)
	rand.Read(about)
	data := map[string]interface{}{"name": "smurfs", "age": 22, "date": time.Unix(1257894000, 0), "about": about}
	body, _ := json.Marshal(data)
	b.StartTimer()
	totalBytes := 0
	sets := 0
	BulkSendor = func(buf *bytes.Buffer) {
		totalBytes += buf.Len()
		sets += 1
		BulkSend(buf)
	}
	for i := 0; i < b.N; i++ {
		IndexBulk("users", "user", strconv.Itoa(i), nil, body)
	}
	log.Printf("Sent %d messages in %d sets totaling %d bytes \n", b.N, sets, totalBytes)
	if BulkErrorCt != 0 {
		b.Fail()
	}
}
