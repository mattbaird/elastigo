package core

import (
	"bytes"
	"encoding/json"
	"github.com/mattbaird/elastigo/api"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	// channel for sending to background indexor
	bulkChannel = make(chan []byte, 100)
	// Max buffer size in bytes before flushing to elasticsearch
	BulkMaxBuffer = 1048576
	// Max number of Docs to hold in buffer before forcing flush
	BulkMaxDocs = 100
	// Max delay before forcing a flush to Elasticearch
	BulkDelaySeconds = 5
	// Keep a running total of errors seen, since it is in the background
	BulkErrorCt uint64
	// We are creating a variable defining the func responsible for sending
	// to allow a mock sendor for test purposes
	BulkSendor func(*bytes.Buffer)
)

// Start up goroutines, channels to start buffering and sending bulk index operations
// The send is a callback incase you want to do a special send, or otherwise see/count etc
// Args
//    @maxConns is the max number of in flight http requests
//    @done is a channel to cause the indexor to stop
//
//   done := make(chan bool)
//   BulkIndexorRun(100, done)
func BulkIndexorRun(maxConns int, done chan bool) {

	go func() {
		bi := NewBulkIndexor(maxConns)
		if BulkSendor == nil {
			BulkSendor = BulkSend
		}
		bi.startHttpSendor()
		bi.startDocChannel()
		bi.startTimer()
		<-done
	}()

}

type BulkIndexor struct {
	sendBuf chan *bytes.Buffer
	buf     *bytes.Buffer
	// Number of documents we have send through so far on this session
	docCt int
	// Max number of http connections in flight at one time
	maxConns int
	// Was the last send induced by time?  or if not, by max docs/size?
	lastSendorByTime bool
	mu               sync.Mutex
}

func NewBulkIndexor(maxConns int) *BulkIndexor {
	b := BulkIndexor{sendBuf: make(chan *bytes.Buffer, maxConns)}
	b.lastSendorByTime = true
	b.buf = new(bytes.Buffer)
	b.maxConns = maxConns
	return &b
}

func (b *BulkIndexor) startHttpSendor() {

	// this sends http requests to elasticsearch it uses maxConns to open up that
	// many goroutines, each of which will synchronously call ElasticSearch
	// in theory, the whole set will cause a backup all the way to IndexBulk if 
	// we have consumed all maxConns
	for i := 0; i < b.maxConns; i++ {
		go func() {
			for {
				buf := <-b.sendBuf
				BulkSendor(buf)
			}
		}()
	}
}

// start a timer for checking back and forcing flush ever BulkDelaySeconds seconds
// even if we haven't hit max messages/size
func (b *BulkIndexor) startTimer() {
	log.Println("Starting timer with delay = ", BulkDelaySeconds)
	ticker := time.NewTicker(time.Second * time.Duration(BulkDelaySeconds))
	go func() {
		for _ = range ticker.C {
			b.mu.Lock()
			// don't send unless last sendor was the time,
			// otherwise an indication of other thresholds being hit
			// where time isn't needed
			if b.buf.Len() > 0 && b.lastSendorByTime {
				b.lastSendorByTime = true
				b.send(b.buf)
			}
			b.mu.Unlock()

		}
	}()
}

func (b *BulkIndexor) startDocChannel() {
	// This goroutine accepts incoming byte arrays from the IndexBulk function and
	// writes to buffer
	go func() {
		for docBytes := range bulkChannel {
			b.mu.Lock()
			b.docCt += 1
			b.buf.Write(docBytes)
			if b.buf.Len() >= BulkMaxBuffer || b.docCt >= BulkMaxDocs {
				b.lastSendorByTime = false
				//log.Printf("Send due to size:  docs=%d  bufsize=%d", b.docCt, b.buf.Len())
				b.send(b.buf)
			}
			b.mu.Unlock()
		}
	}()
}

func (b *BulkIndexor) send(buf *bytes.Buffer) {
	//b2 := *b.buf
	b.sendBuf <- buf
	b.buf = new(bytes.Buffer)
	b.docCt = 0
}

// This does the actual send of a buffer, which has already been formatted
// into bytes of ES formatted bulk data
func BulkSend(buf *bytes.Buffer) {
	_, err := api.DoCommand("POST", "/_bulk", buf)
	if err != nil {
		log.Println(err)
		BulkErrorCt += 1
	}
}

// The index bulk API adds or updates a typed JSON document to a specific index, making it searchable. 
// it operates by buffering requests, and ocassionally flushing to elasticsearch
// http://www.elasticsearch.org/guide/reference/api/bulk.html
func IndexBulk(index string, _type string, id string, date *time.Time, data interface{}) error {
	//{ "index" : { "_index" : "test", "_type" : "type1", "_id" : "1" } }
	buf := bytes.Buffer{}
	buf.WriteString(`{"index":{"_index":"`)
	buf.WriteString(index)
	buf.WriteString(`","_type":"`)
	buf.WriteString(_type)
	buf.WriteString(`","_id":"`)
	buf.WriteString(id)
	if date != nil {
		buf.WriteString(`","_timestamp":"`)
		buf.WriteString(strconv.FormatInt(date.UnixNano()/1e6, 10))
	}
	buf.WriteString(`"}}`)
	buf.WriteByte('\n')
	switch v := data.(type) {
	case *bytes.Buffer:
		io.Copy(&buf, v)
	case []byte:
		buf.Write(v)
	case string:
		buf.WriteString(v)
	default:
		body, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			log.Println("Json data error ", data)
			return jsonErr
		}
		buf.Write(body)
	}
	buf.WriteByte('\n')
	bulkChannel <- buf.Bytes()
	return nil
}
