package job

import (
  "os"
  "errors"
  "math/rand"
  "time"
  "lib/message"
  "encoding/binary"
  "encoding/hex"
  "lib/job/jobmsg"
)

type Job struct {
  ID string
  rate int
  syslogFacility int
  syslogSeverity int
  fileName string
  fileHandle os.File
  jobChan  chan
  conn net.Conn
}

func NewJob(r int, sf int, ss int, file string, jc chan) (Job,error) {
  j := Job{
    rate:r,
    syslogFacility: sf,
    syslogSeverity: ss,
    fileName: file,
    jobChan: jc
  }
  err := j.openFile()
  j.genID()
  return j,err
}

func (j *Job) genID () {
  //generate new random ID
  rand.Seed( time.Now().UTC().UnixNano())
  randNum := rand.Intn(18446744073709551615 - 1) + 1
  bytes := [8]byte{}
  binary.LittleEndian.PutUint64(bytes[:],randNum)
  j.ID = hex.EncodeToString(bytes[:s])
}

func (j *Job) openFile() error {
  file, err := os.Open(j.fileName)
  if err != nil {
    //handle file error
    //report back that the file cant be opened and why
    err = errors.New("Unable to open file")
  }
  j.fileHandle = file
  return err
}

func (j *Job) closeFile(){
  j.fh.Close()
}

func (j *Job) Start(){
  //start the message sender
  //set it as a go routine with a channel to be controlled
  ticker := time.NewTicker(time.Second * 1)
  fileRead := bufio.NewReader(j.fileHandle)

  for _ = range ticker.C {
    for i := 0; i < *rate; i++ {
      lineBuffer, _, err := fileRead.ReadLine()
      if err != nil {
        //log.Println(err)
        j.openFile()
        fileRead = bufio.NewReader(j.fileHandle)
        lineBuffer, _, err = fileRead.ReadLine()
      }
      msg := message.NewMessage(sourceHost, syslogFacility, syslogSeverity)
      msg.AddToMessage(string(lineBuffer))
      msg.Send(j.conn)
    }
    //Check for close requests
    select {
      case jobMsg := <- j.jobChan:
        if jobMsg.ID == j.ID {
          //message is for me handle its request
          if jobMsg.action == jobmsg.stop {
            //exit go task
            //send message job ID is stopping
            j.closeFile()
            return
          } else if jobMsg.action == jobmsg.start {
            //????
          } else if jobMsg.action == jobmsg.restart{
            //restart task
            //this leaves the job configured but it forces it to restart
          }
        }
    }
  }
}
