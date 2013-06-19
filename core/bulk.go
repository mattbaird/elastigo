package core

import (
	"bytes"
	"encoding/json"
	u "github.com/araddon/gou"
	"github.com/mattbaird/elastigo/api"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	// Max buffer size in bytes before flushing to elasticsearch
	BulkMaxBuffer = 1048576
	// Max number of Docs to hold in buffer before forcing flush
	BulkMaxDocs = 100
	// Max delay before forcing a flush to Elasticearch
	BulkDelaySeconds = 5
	// Keep a running total of errors seen, since it is in the background
	BulkErrorCt uint64

	// There is one Global Bulk Indexor for convenience
	bulkIndexor *BulkIndexor
)

type ErrorBuffer struct {
	Err error
	Buf *bytes.Buffer
}

// There is one global bulk indexor available for convenience so the IndexBulk() function can be called.
// However, the recommended usage is create your own BulkIndexor to allow for multiple seperate elasticsearch
// servers/host connections.
//    @maxConns is the max number of in flight http requests
//    @done is a channel to cause the indexor to stop
//
//   done := make(chan bool)
//   BulkIndexorGlobalRun(100, done)
func BulkIndexorGlobalRun(maxConns int, done chan bool) {
	if bulkIndexor == nil {
		bulkIndexor = NewBulkIndexor(maxConns)
		bulkIndexor.Run(done)
	}
}

// A bulk indexor creates goroutines, and channels for connecting and sending data
// to elasticsearch in bulk, using buffers.
type BulkIndexor struct {

	// We are creating a variable defining the func responsible for sending
	// to allow a mock sendor for test purposes
	BulkSendor func(*bytes.Buffer) error

	// If we encounter an error in sending, we are going to retry for this long
	// before returning an error
	// if 0 it will not retry
	RetryForSeconds int

	// channel for getting errors
	ErrorChannel chan *ErrorBuffer

	// channel for sending to background indexor
	bulkChannel chan []byte

	// shutdown channel
	shutdownChan chan bool

	// buffers
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
	b.bulkChannel = make(chan []byte, 100)
	return &b
}

// A bulk indexor with more control over error handling
//    @maxConns is the max number of in flight http requests
//    @retrySeconds is # of seconds to wait before retrying falied requests
//
//   done := make(chan bool)
//   BulkIndexorGlobalRun(100, done)
func NewBulkIndexorErrors(maxConns, retrySeconds int) *BulkIndexor {
	b := BulkIndexor{sendBuf: make(chan *bytes.Buffer, maxConns)}
	b.lastSendorByTime = true
	b.buf = new(bytes.Buffer)
	b.maxConns = maxConns
	b.RetryForSeconds = retrySeconds
	b.bulkChannel = make(chan []byte, 100)
	b.ErrorChannel = make(chan *ErrorBuffer, 20)
	return &b
}

// Starts this bulk Indexor running, this Run opens a go routine so is
// Non blocking
func (b *BulkIndexor) Run(done chan bool) {

	go func() {
		if b.BulkSendor == nil {
			b.BulkSendor = BulkSend
		}
		b.shutdownChan = done
		b.startHttpSendor()
		b.startDocChannel()
		b.startTimer()
		<-b.shutdownChan
	}()
}

// Flush all current documents to ElasticSearch
func (b *BulkIndexor) Flush() {
	b.mu.Lock()
	if b.docCt > 0 {
		b.send(b.buf)
	}
	b.mu.Unlock()
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
				err := b.BulkSendor(buf)

				// Perhaps a b.FailureStrategy(err)  ??  with different types of strategies
				//  1.  Retry, then panic
				//  2.  Retry then return error and let runner decide
				//  3.  Retry, then log to disk?   retry later?
				if err != nil {
					if b.RetryForSeconds > 0 {
						time.Sleep(time.Second * time.Duration(b.RetryForSeconds))
						err = b.BulkSendor(buf)
						if err == nil {
							continue
						}
					}
					if b.ErrorChannel != nil {
						log.Println(err)
						b.ErrorChannel <- &ErrorBuffer{err, buf}
					}
				}
			}
		}()
	}
}

// start a timer for checking back and forcing flush ever BulkDelaySeconds seconds
// even if we haven't hit max messages/size
func (b *BulkIndexor) startTimer() {
	u.Debug("Starting Bulk timer with delay = ", BulkDelaySeconds)
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
		for docBytes := range b.bulkChannel {
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

// The index bulk API adds or updates a typed JSON document to a specific index, making it searchable.
// it operates by buffering requests, and ocassionally flushing to elasticsearch
// http://www.elasticsearch.org/guide/reference/api/bulk.html
func (b *BulkIndexor) Index(index string, _type string, id string, date *time.Time, data interface{}) error {
	//{ "index" : { "_index" : "test", "_type" : "type1", "_id" : "1" } }
	by, err := IndexBulkBytes(index, _type, id, date, data)
	if err != nil {
		u.Error(err)
		return err
	}
	b.bulkChannel <- by
	return nil
}

// This does the actual send of a buffer, which has already been formatted
// into bytes of ES formatted bulk data
func BulkSend(buf *bytes.Buffer) error {
	_, err := api.DoCommand("POST", "/_bulk", buf)
	if err != nil {
		log.Println(err)
		BulkErrorCt += 1
		return err
	}
	return nil
}

// Given a set of arguments for index, type, id, data create a set of bytes that is formatted for bulkd index
// http://www.elasticsearch.org/guide/reference/api/bulk.html
func IndexBulkBytes(index string, _type string, id string, date *time.Time, data interface{}) ([]byte, error) {
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
			return nil, jsonErr
		}
		buf.Write(body)
	}
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

// The index bulk API adds or updates a typed JSON document to a specific index, making it searchable.
// it operates by buffering requests, and ocassionally flushing to elasticsearch
// http://www.elasticsearch.org/guide/reference/api/bulk.html
func IndexBulk(index string, _type string, id string, date *time.Time, data interface{}) error {
	//{ "index" : { "_index" : "test", "_type" : "type1", "_id" : "1" } }
	if bulkIndexor == nil {
		panic("Must have Global Bulk Indexor to use this Func")
	}
	by, err := IndexBulkBytes(index, _type, id, date, data)
	if err != nil {
		return err
	}
	bulkIndexor.bulkChannel <- by
	return nil
}
