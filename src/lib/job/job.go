package job

import (
  "os"
  "errors"
  "net"
  "math/rand"
  "time"
  "bufio"
  "log"
  "lib/message"
  "lib/job/jobmsg"
  "encoding/binary"
  "encoding/hex"
)

type Job struct {
  ID string
  rate *int
  syslogFacility *int
  syslogSeverity *int
  sourceHost *string
  fileName *string
  fileHandle *os.File
  ctrlChannel chan jobmsg.JobMsg
  conn net.Conn
}

func NewJob(dip *string, dport *string,r *int, sf *int, ss *int, sh *string, file *string, cc chan jobmsg.JobMsg) (Job,error) {
  destAddr, err := net.ResolveUDPAddr("udp", *dip+":"+*dport)
  con, err := net.DialUDP("udp", nil, destAddr)
  j := Job{
    rate:r,
    sourceHost: sh,
    syslogFacility: sf,
    syslogSeverity: ss,
    fileName: file,
    ctrlChannel: cc,
    conn:con,
  }
  err = j.openFile()
  return j,err
}

func (j *Job) GenID () string {
  //generate new random ID
  rand.Seed( time.Now().UTC().UnixNano())
  randNum := rand.Uint32() + 1
  bytes := [4]byte{}
  binary.LittleEndian.PutUint32(bytes[:],randNum)
  j.ID = hex.EncodeToString(bytes[:])
  return j.ID
}

func (j *Job) openFile() error {
  file, err := os.Open("/home/rcameron/code/gogoLogs/src/" + *j.fileName)
  if err != nil {
    //handle file error
    //report back that the file cant be opened and why
    log.Println(err)
    err = errors.New("Unable to open file")
  }
  j.fileHandle = file
  return err
}

func (j *Job) closeFile(){
  j.fileHandle.Close()
}

func (j *Job) Start(){
  //start the message sender
  //set it as a go routine with a channel to be controlled
  ticker := time.NewTicker(time.Second * 1)
  fileRead := bufio.NewReader(j.fileHandle)

  for _ = range ticker.C {
    log.Println("tick",*j.rate)
    for i := 0; i < *j.rate; i++ {
      lineBuffer, _, err := fileRead.ReadLine()
      if err != nil {
        log.Println(err)
        j.openFile()
        fileRead = bufio.NewReader(j.fileHandle)
        lineBuffer, _, err = fileRead.ReadLine()
      }
      msg := message.NewMessage(j.sourceHost, j.syslogFacility, j.syslogSeverity)
      msg.AddToMessage(string(lineBuffer))
      log.Println(msg)
      msg.Send(j.conn)
    }
    //Check for close requests
  }
}
