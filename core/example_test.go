package core_test

import (
	"fmt"
	"github.com/mattbaird/elastigo/core"
	"strconv"
	"time"
)

// The simplest usage of background bulk indexing
func ExampleBulkIndexor_simple() {
	indexor := core.NewBulkIndexorErrors(10, 60)
	done := make(chan bool)
	indexor.Run(done)

	indexor.Index("twitter", "user", "1", "", nil, `{"name":"bob"}`)

	<-done // wait forever
}

// The simplest usage of background bulk indexing with error channel
func ExampleBulkIndexor_errorchannel() {
	indexor := core.NewBulkIndexorErrors(10, 60)
	done := make(chan bool)
	indexor.Run(done)

	go func() {
		for errBuf := range indexor.ErrorChannel {
			// just blissfully print errors forever
			fmt.Println(errBuf.Err)
		}
	}()
	for i := 0; i < 20; i++ {
		indexor.Index("twitter", "user", strconv.Itoa(i), "", nil, `{"name":"bob"}`)
	}
	<-done
}

// The simplest usage of background bulk indexing with error channel
func ExampleBulkIndexor_errorsmarter() {
	indexor := core.NewBulkIndexorErrors(10, 60)
	done := make(chan bool)
	indexor.Run(done)

	errorCt := 0 // use sync.atomic or something if you need
	timer := time.NewTicker(time.Minute * 3)
	go func() {
		for {
			select {
			case _ = <-timer.C:
				if errorCt < 2 {
					errorCt = 0
				}
			case _ = <-done:
				return
			}
		}
	}()

	go func() {
		for errBuf := range indexor.ErrorChannel {
			errorCt++
			fmt.Println(errBuf.Err)
			// log to disk?  db?   ????  Panic
		}
	}()
	for i := 0; i < 20; i++ {
		indexor.Index("twitter", "user", strconv.Itoa(i), "", nil, `{"name":"bob"}`)
	}
	<-done
}

// The inspecting the response
func ExampleBulkIndexor_responses() {
	indexor := core.NewBulkIndexor(10)
	// Create a custom Sendor Func, to allow inspection of response/error
	indexor.BulkSendor = func(buf *bytes.Buffer) error {
		// @buf is the buffer of docs about to be written
		respJson, err := api.DoCommand("POST", "/_bulk", buf)
		if err != nil {
			// handle it better than this
			fmt.Println(string(respJson))
		}
		return err
	}
	done := make(chan bool)
	indexor.Run(done)

	for i := 0; i < 20; i++ {
		indexor.Index("twitter", "user", strconv.Itoa(i), "", nil, `{"name":"bob"}`)
	}
	<-done
}
