package main

import (
	"lib/webui/handlers"
	"lib/message"
	"lib/job/jobmsg"
	"lib/stats"
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"time"
)

/*Const used in the app*/
const (
	msgBufSize = 4096
)

/*Arguments possible to send to tool*/
var destinationIP = flag.String("d", "127.0.0.1", "Specify destination IP (default: 127.0.0.1)")
var sourceHost = flag.String("s", "127.0.0.1", "Specify source hostname/IP for syslog header (default: 127.0.0.1)")
var destinationPort = flag.String("p", "514", "Sepecify port (default: 514)")
var protocol = flag.String("P", "UDP", "Specify protocol to send data over (default: UDP)")
var fileName = flag.String("f", "", "Specify file to read from (default: None)")
var rate = flag.Int("r", 5, "Specify rate (default: 5/s)")
var syslogFacility = flag.Int("F", 0, "Specify syslog priority value (default:0)")
var syslogSeverity = flag.Int("S", 0, "Specify syslog priority value (default:0)")
var nonStop = flag.Bool("C", false, "Specify if the file should be continously read from (default: false)")

//WebUI
var enableWebUI = flag.Bool("w", false, "Enable WebIU for log sender (default: false)")
var webUIPort = flag.Int("wP",8080, "Specify the port to listen on (default: 8080) (Used with -w)")
var dirName = flag.String("wD", "", "Specify the directory to serve logs from (default: None, used with -w)")
//create channels to handle listening to messages

/*handleMessages a go routine to handle read files */
func handleMessages(conn net.Conn, rate *int, sendChannel chan message.Message) {
	//ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case msg := <-sendChannel:
			msg.Send(conn)
		}
	}
}

/*Main entry point*/
func main() {
	/* initialize channels */
	sendChannel := make(chan message.Message, msgBufSize)
	/* Parse command line flags */
	flag.Parse()

	if *enableWebUI != true {
		//stand alone mode
		file, err := os.Open(*fileName)
		if err != nil {
			log.Fatal(err)
		}

		//create UDP connection
		//allow user to specify TCP or UDP
		destAddr, err := net.ResolveUDPAddr("udp", *destinationIP+":"+*destinationPort)
		con, err := net.DialUDP("udp", nil, destAddr)

		fileRead := bufio.NewReader(file)

		go handleMessages(con, rate, sendChannel)

		for {
			ticker := time.NewTicker(time.Second * 1)
			for _ = range ticker.C {
				for i := 0; i < *rate; i++ {
					lineBuffer, _, err := fileRead.ReadLine()
					if err != nil {
						//log.Println(err)
						if *nonStop {
							file.Close()
							file, _ := os.Open(*fileName)
							fileRead = bufio.NewReader(file)
							lineBuffer, _, err = fileRead.ReadLine()
						} else {
						}
					}
					msg := message.NewMessage(sourceHost, syslogFacility, syslogSeverity)
					msg.AddToMessage(string(lineBuffer))
					sendChannel <- msg
				}
			}
		}
	} else if *enableWebUI == true {
		//create channels
		cc := make(chan jobmsg.JobMsg,msgBufSize)
		sc := make(chan stats.Stats,msgBufSize)
		handy := handlers.NewHandler(cc,sc,*webUIPort,*dirName)
		handy.Start()
	}
}
