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

func calcSyslogPriority(f *int, s *int ) int {
	return (*f * 8 ) + *s
}

func main() {
	ticker := time.NewTicker(time.Second * 1)
	flag.Parse()
	syslogPriority := calcSyslogPriority(syslogFacility,syslogSeverity)

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	//create UDP connection
	//allow user to specify TCP or UDP
	destAddr, err := net.ResolveUDPAddr("udp",*destinationIP + ":" + *destinationPort)
	con, err := net.DialUDP("udp",nil,destAddr)

	scanner := bufio.NewScanner(file)

	for _ = range ticker.C {
		for i := 0; i < *rate; i ++ {
			EOFmet := scanner.Scan()
			if EOFmet != true {
				if *nonStop != true {
					os.Exit(0)
				}
			}
			//Write format to syslog standard
			var message []string
			//add priority to string
			message = append(message,"<" + strconv.Itoa(syslogPriority) + ">")
			//Add standard syslog timestamp RFC3339 format
			message = append(message,time.Now().Format(time.RFC3339))
			//Add in source IP/hostname
			message = append(message,*sourceHost)
			//add in payload
			message = append(message,scanner.Text())

			finalMessage := strings.Join(message," ")
			con.Write([]byte(finalMessage))
			log.Println(finalMessage)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Scanner err")
		log.Fatal(err)
	}
}
