package search

import (
	"flag"
	"github.com/araddon/gou"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"log"
	"os"
)

var (
	_                 = log.Ldate
	hasStartedTesting bool
	eshost            *string = flag.String("host", "localhost", "Elasticsearch Server Host Address")
	logLevel          *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

/*

usage:

	test -v -host eshost 

*/

func init() {
	InitTests(false)
	core.DebugRequests = true
}

func InitTests(startIndexor bool) {
	if !hasStartedTesting {
		flag.Parse()
		hasStartedTesting = true
		gou.SetLogger(log.New(os.Stderr, "", log.Ltime|log.Lshortfile), *logLevel)
		log.SetFlags(log.Ltime | log.Lshortfile)
		api.Domain = *eshost
	}
}
