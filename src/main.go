package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"net"
	"time"
)

var destinationIP = flag.String("d", "127.0.0.1", "Specify destination IP (default: 514)")
var destinationPort = flag.String("p", "514", "Sepecify port (default: 53)")
var fileName = flag.String("f", "", "Specify file to read from")
var rate = flag.Int("r",5,"Specify rate (default: 5/s)")

func main() {
	ticker := time.NewTicker(time.Second * 1)
	flag.Parse()
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	destAddr, err := net.ResolveUDPAddr("udp",*destinationIP + ":" + *destinationPort)
	con, err := net.DialUDP("udp",nil,destAddr)

	scanner := bufio.NewScanner(file)

	for _ = range ticker.C {
		for i := 0; i < *rate; i ++ {
			scanner.Scan()
			con.Write([]byte(scanner.Text()))
			log.Println(scanner.Pos())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Scanner err")
		log.Fatal(err)
	}
}
