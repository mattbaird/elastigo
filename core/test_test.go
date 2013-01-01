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

This loads test data into elasticsearch for testing
IT should not get compiled into package (hopefully, because the test_load.go name?)

usage:

	test -v -host eshost -loaddata 

https://api.github.com/events

	{
		"created_at":"2012-04-11T15:09:50Z",
		"payload":{
			"ref":"refs/heads/master",
			"commits":[{"message":"almost done","sha":"c680f05d6bcce30b72a98d001ec5717c6dff26b4","distinct":true,"url":"https://api.github.com/repos/hydna/hydna-ruby/commits/c680f05d6bcce30b72a98d001ec5717c6dff26b4","author":{"name":"Isak Wiström","email":"isak.wistrom@gmail.com"}}],
			"push_id":72207903,
			"size":1,
			"head":"c680f05d6bcce30b72a98d001ec5717c6dff26b4"
		},
		"repo":{
			"url":"https://api.github.dev/repos/hydna/hydna-ruby","name":"hydna/hydna-ruby","id":3993348
		},
		"type":"PushEvent",
		"public":true,
		"org":{
			"gravatar_id":"b91d9862d1e1bc8c1089bf4bf93dd51f",
			"url":"https://api.github.dev/orgs/hydna","login":"hydna","avatar_url":"https://secure.gravatar.com/avatar/b91d9862d1e1bc8c1089bf4bf93dd51f?d=http://github.dev%2Fimages%2Fgravatars%2Fgravatar-org-420.png","id":194557},
		"actor":{
				"gravatar_id":"216d18469ecac1eda15368471754d2ab",
				"url":"https://api.github.dev/users/Skaggivara",
				"login":"Skaggivara",
				"avatar_url":"https://secure.gravatar.com/avatar/216d18469ecac1eda15368471754d2ab?d=http://github.dev%2Fimages%2Fgravatars%2Fgravatar-user-420.png",
				"id":194565
		},
		"id":"1540149912"
	}

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

		flag.Parse()
		log.SetFlags(log.Ltime | log.Lshortfile)
		api.Domain = *eshost
	}
	if startIndexor && !bulkStarted {
		BulkDelaySeconds = 1
		bulkStarted = true
		log.Println("start bulk indexor")
		BulkIndexorRun(100, make(chan bool))
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
	BulkSendor = func(buf *bytes.Buffer) {
		log.Printf("Sent %d bytes total %d docs sent", buf.Len(), docCt)
		BulkSend(buf)
	}
	resp, err := http.Get("http://data.githubarchive.org/2012-12-10-15.json.gz")
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
			break
		}
		if err := json.Unmarshal(line, &ge); err == nil {
			// obviously there is some chance of collision here so only useful for testing
			// plus, i don't even know if the url is unique?
			id := strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(ge.Url))), 10)
			IndexBulk("github", ge.Type, id, &ge.Created, line)
			docCt++
			//log.Println(string(line))
			//os.Exit(1)
		} else {
			log.Println("ERROR? ", string(line))
		}

	}
}
