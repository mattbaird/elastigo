package search

import (
	"flag"
	"github.com/araddon/gou"
	"github.com/meanpath/elastigo/api"
	"github.com/meanpath/elastigo/core"
	"log"
	"os"
	//"testing"
)

var (
	_                 = log.Ldate
	hasStartedTesting bool
	eshost            *string = flag.String("host", "localhost", "Elasticsearch Server Host Address")
	logLevel          *string = flag.String("logging", "info", "Which log level: [debug,info,warn,error,fatal]")
)

/*

usage:

	test -v -host eshost

*/

func init() {
	InitTests(false)
	if *logLevel == "debug" {
		//*logLevel = "debug"
		core.DebugRequests = true
	}
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
