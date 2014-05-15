package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"net"
	"strconv"
	"strings"
	"time"
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

 type Message struct {
	message []string
	sourceHost *string
	syslogFacility *int
	syslogSeverity *int
	syslogPriority int
}

func (m *Message) setMessageTime() {
	m.message = append(m.message,time.Now().Format(time.RFC3339))
}

func (m *Message) setSyslogPriority() {
	m.message = append(m.message,"<" + strconv.Itoa(m.syslogPriority) + ">")
}

func (m *Message) setSrcHost() {
	m.message = append(m.message,*m.sourceHost)
}

func (m *Message) AddToMessage(s string) {
	m.message = append(m.message,s)
}

func (m *Message) calcSyslogPriority(f *int, s *int ) {
	m.syslogPriority = (*f * 8 ) + *s
}

func (m *Message) send(con net.Conn) {
	finalMessage := strings.Join(m.message," ")
	con.Write([]byte(finalMessage))
	//log.Println(finalMessage)
}

func NewMessage(srcHost *string, f *int, s *int) Message {
	msg := Message{sourceHost:srcHost,syslogFacility:f, syslogSeverity:s}
	msg.calcSyslogPriority(f,s)
	msg.setSyslogPriority()
	msg.setMessageTime()
	msg.setSrcHost()
	return msg
}

/*handleMessages a go routine to handle read files */
func handleMessages(conn net.Conn, rate *int, sendChannel chan Message) {
	//ticker := time.NewTicker(time.Second * 1)
	for {
		select {
			case msg := <-sendChannel:
				msg.send(conn)
		}
	}
}


func main() {
	/* initialize channels */
	sendChannel := make(chan Message, msgBufSize)
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
				msg := NewMessage(sourceHost,syslogFacility,syslogSeverity)
				msg.AddToMessage(string(lineBuffer))
				sendChannel <- msg
			}
		}
	}
}
