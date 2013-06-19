package core_test

import (
	"fmt"
	"github.com/mattbaird/elastigo/core"
	"time"
)

// The simplest usage of background bulk indexing
func ExampleBulkIndexor_simple() {
	indexor := core.NewBulkIndexorErrors(10, 60)
	done := make(chan bool)
	indexor.Run(done)

	<-done // wait forever
}

// The simplest usage of background bulk indexing with error channel
func ExampleBulkIndexor_errorchannel() {
	indexor := core.NewBulkIndexorErrors(10, 60)
	done := make(chan bool)
	indexor.Run(done)

	for errBuf := range indexor.ErrorChannel {
		// just blissfully print errors forever
		fmt.Println(errBuf.Err)
	}
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
	for errBuf := range indexor.ErrorChannel {
		errorCt++
		fmt.Println(errBuf.Err)
		// log to disk?  db?   ????  Panic
	}
}
