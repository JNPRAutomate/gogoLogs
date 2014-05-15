package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"net"
	"time"
	"./lib/message"
)

const (
	msgBufSize = 4096
)

var destinationIP = flag.String("d", "127.0.0.1", "Specify destination IP (default: 127.0.0.1)")
var sourceHost = flag.String("s","127.0.0.1", "Specify source hostname/IP for syslog header (default: 127.0.0.1)")
var destinationPort = flag.String("p", "514", "Sepecify port (default: 514)")
var protocol = flag.String("P","UDP","Specify protocol to send data over (default: UDP)")
var fileName = flag.String("f", "", "Specify file to read from (default: None)")
var rate = flag.Int("r",5,"Specify rate (default: 5/s)")
var syslogFacility = flag.Int("F",0,"Specify syslog priority value (default:0)")
var syslogSeverity = flag.Int("S",0,"Specify syslog priority value (default:0)")
var nonStop = flag.Bool("C",false,"Specify if the file should be continously read from (default: false)")
//WebUI
var enableWebUI = flag.Bool("w",false,"Enable WebIU for log sender (default: false) (To be added)")

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


func main() {
	/* initialize channels */
	sendChannel := make(chan message.Message, msgBufSize)
	/* Parse command line flags */
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	//create UDP connection
	//allow user to specify TCP or UDP
	destAddr, err := net.ResolveUDPAddr("udp",*destinationIP + ":" + *destinationPort)
	con, err := net.DialUDP("udp",nil,destAddr)

	fileRead := bufio.NewReader(file)

	go handleMessages(con,rate,sendChannel)

	for {
		ticker := time.NewTicker(time.Second * 1)
		for _ = range ticker.C {
			for i := 0; i < *rate; i ++ {
				lineBuffer, _, err := fileRead.ReadLine()
				if err != nil {
					//log.Println(err)
					if *nonStop  {
						file, _ := os.Open(*fileName)
						fileRead = bufio.NewReader(file)
						lineBuffer, _, err = fileRead.ReadLine()
					} else {
					}
				}
				msg := message.NewMessage(sourceHost,syslogFacility,syslogSeverity)
				msg.AddToMessage(string(lineBuffer))
				sendChannel <- msg
			}
		}
	}
}
