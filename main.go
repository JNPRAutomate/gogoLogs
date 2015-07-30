package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"./lib/job/jobmsg"
	"./lib/stats"
	"./lib/webui/handlers"
)

/*Const used in the app*/
const (
	msgBufSize = 4096
)

//Flags
var webUIPort = flag.Int("wP", 8080, "Optional: Specify the port to listen on (default: 8080) (Used with -w)")
var dirName = flag.String("wD", "", "Required: Specify the directory to serve logs from (default: None, used with -w)")
var destHosts = flag.String("d", "", "Optional: Specify the destination hosts to prepopulate this field in the UI. Use comma to provide multiple values. (Example: 1.2.3.4 or 1.2.3.4,2.3.4.5)")
var sourceNames = flag.String("s", "", "Optional: Specify the source hosts to prepopulate this field in the UI. Use comma to provide multiple values. (Example: foo or foo,bar)")

//create channels to handle listening to messages

/*Main entry point*/
func main() {
	var destSlice []string
	var sourceNamesSlice []string

	/* Parse command line flags */
	flag.Parse()

	if *dirName != "" {
		dirTest, err := os.Open(*dirName)
		if err != nil {
			log.Fatalln(err)
		}
		dirTest.Close()

		if *destHosts != "" {
			destSlice = strings.Split(*destHosts, ",")
			log.Printf("Setting destination hosts to: %s", destSlice)
		} else {
			destSlice = make([]string, 0)
		}

		if *sourceNames != "" {
			sourceNamesSlice = strings.Split(*sourceNames, ",")
			log.Printf("Setting source hosts to: %s", sourceNamesSlice)
		} else {
			sourceNamesSlice = make([]string, 0)
		}

		//create channels
		cc := make(chan jobmsg.JobMsg, msgBufSize)
		sc := make(chan stats.Stats, msgBufSize)
		handy := handlers.NewHandler(cc, sc, *webUIPort, *dirName, destSlice, sourceNamesSlice)
		handy.Start()
	} else {
		flag.PrintDefaults()
	}
}
