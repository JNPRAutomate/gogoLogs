package message

import (
  "net"
  "strconv"
  "strings"
  "time"
)

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

func (m *Message) Send(con net.Conn) {
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
