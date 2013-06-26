package core

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"github.com/mattbaird/elastigo/api"
	"hash/crc32"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

/*

usage:

	test -v -host eshost -loaddata

*/

var (
	_                 = os.ModeDir
	bulkStarted       bool
	hasLoadedData     bool
	hasStartedTesting bool
	eshost            *string = flag.String("host", "localhost", "Elasticsearch Server Host Address")
	loadData          *bool   = flag.Bool("loaddata", false, "This loads a bunch of test data into elasticsearch for testing")
)

func init() {

}
func InitTests(startIndexor bool) {
	if !hasStartedTesting {
		flag.Parse()
		hasStartedTesting = true
		log.SetFlags(log.Ltime | log.Lshortfile)
		api.Domain = *eshost
	}
	if startIndexor && !bulkStarted {
		BulkDelaySeconds = 1
		bulkStarted = true
		log.Println("start bulk indexor")
		BulkIndexorGlobalRun(100, make(chan bool))
		if *loadData && !hasLoadedData {
			log.Println("load test data ")
			hasLoadedData = true
			LoadTestData()
		}
	}
}

// dumb simple assert for testing, printing
//    Assert(len(items) == 9, t, "Should be 9 but was %d", len(items))
func Assert(is bool, t *testing.T, format string, args ...interface{}) {
	if is == false {
		log.Printf(format, args...)
		t.Fail()
	}
}

// Wait for condition (defined by func) to be true, a utility to create a ticker
// checking every 100 ms to see if something (the supplied check func) is done
//
//   WaitFor(func() bool {
//      return ctr.Ct == 0
//   },10)
//
// @timeout (in seconds) is the last arg
func WaitFor(check func() bool, timeoutSecs int) {
	timer := time.NewTicker(100 * time.Millisecond)
	tryct := 0
	for _ = range timer.C {
		if check() {
			timer.Stop()
			break
		}
		if tryct >= timeoutSecs*10 {
			timer.Stop()
			break
		}
		tryct++
	}
}

func TestFake(t *testing.T) {

}

type GithubEvent struct {
	Url     string
	Created time.Time `json:"created_at"`
	Type    string
}

// This loads test data from github archives (~6700 docs)
func LoadTestData() {
	docCt := 0
	indexor := NewBulkIndexor(20)
	indexor.BulkSendor = func(buf *bytes.Buffer) error {
		log.Printf("Sent %d bytes total %d docs sent", buf.Len(), docCt)
		return BulkSend(buf)
	}
	resp, err := http.Get("http://data.githubarchive.org/2012-12-10-15.json.gz")
	if err != nil || resp == nil {
		panic("Could not download data")
	}
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}
	gzReader, err := gzip.NewReader(resp.Body)
	defer gzReader.Close()
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(gzReader)
	var ge GithubEvent
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			indexor.Flush()
			break
		}
		if err := json.Unmarshal(line, &ge); err == nil {
			// obviously there is some chance of collision here so only useful for testing
			// plus, i don't even know if the url is unique?
			id := strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(ge.Url))), 10)
			indexor.Index("github", ge.Type, id, "", &ge.Created, line)
			docCt++
			//log.Println(string(line))
			//os.Exit(1)
		} else {
			log.Println("ERROR? ", string(line))
		}

	}
	// lets wait a bit to ensure that elasticsearch finishes?
	time.Sleep(time.Second * 5)
}
