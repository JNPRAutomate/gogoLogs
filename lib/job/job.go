package job

import (
  "os"
  "errors"
  "net"
  "math/rand"
  "time"
  "bufio"
  "strconv"
  "lib/message"
  "lib/job/jobmsg"
  "lib/stats/statsmsg"
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
  fileSize int64
  fileHandle *os.File
  StatsChannel chan statsmsg.StatsMsg
  CtrlChannel chan jobmsg.JobMsg
  conn net.Conn
  jTimer JobTimer
}

type JobTimer struct {
  MaxCount uint
  CurrentCount uint
}

func NewJob(dip *string, dport *string,r *int, sf *int, ss *int, sh *string, file *string, cc chan jobmsg.JobMsg, mc uint) (Job,error) {
  destAddr, err := net.ResolveUDPAddr("udp", *dip+":"+*dport)
  con, err := net.DialUDP("udp", nil, destAddr)
  j := Job{
    rate:r,
    sourceHost: sh,
    syslogFacility: sf,
    syslogSeverity: ss,
    fileName: file,
    CtrlChannel: cc,
    conn:con,
    jTimer: JobTimer{MaxCount:mc,CurrentCount:0},
  }
  err = j.openFile()
  return j,err
}

func (j *Job) GenID () string {
  //generate new random ID
  rand.Seed( time.Now().UTC().UnixNano())
  randNum := rand.Uint32() + 1
  bytes := [4]byte{}
  binary.BigEndian.PutUint32(bytes[:],randNum)
  j.ID = hex.EncodeToString(bytes[:])
  return j.ID
}

func (j *Job) SetID(id int) {
  j.ID = strconv.Itoa(id)
}

func (j *Job) openFile() error {
  var file *os.File
  var err error
  var fileStats os.FileInfo
  file, err = os.Open(*j.fileName)
  if err != nil {
    //handle file error
    //report back that the file cant be opened and why
    //log.Println(err)
    err = errors.New("Unable to open file")
    return err
  }
  j.fileHandle = file

  fileStats, err = file.Stat();
  if err != nil {
    err = errors.New("Unable to stat file")
    return err
  }
  j.fileSize = fileStats.Size()

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
  var sendRate uint
  var totalSent uint
  sendRate = 0
  totalSent = 0

  for _ = range ticker.C {
    select {
      //Check for close requests
      case newJobMsg := <-j.CtrlChannel:
        if newJobMsg.Action == jobmsg.Stop {
          //exiting
          return
        }
      //continue job
      default:
        for i := 0; i < *j.rate; i++ {
          lineBuffer, _, err := fileRead.ReadLine()
          if err != nil {
            j.fileHandle.Close()
            //log.Println(err)
            j.openFile()
            fileRead = bufio.NewReader(j.fileHandle)
            lineBuffer, _, err = fileRead.ReadLine()
          }
          msg := message.NewMessage(j.sourceHost, j.syslogFacility, j.syslogSeverity)
          msg.AddToMessage(string(lineBuffer))
          msg.Send(j.conn)
          totalSent = totalSent + 1
          sendRate = sendRate + 1
        }
        j.jTimer.CurrentCount = j.jTimer.CurrentCount + 1;
        //send stats here
        // sentRate, totalSent, jobID
        //reset sent rate
        j.StatsChannel <- statsmsg.StatsMsg{ID:j.ID,TotalSent: totalSent, SendRate: sendRate}
        sendRate = 0
        //Check timer
        if j.jTimer.MaxCount > 0 {
          if j.jTimer.CurrentCount == j.jTimer.MaxCount {
            //exit its max life
            //message job manager
            j.StatsChannel <- statsmsg.StatsMsg{ID:j.ID,TotalSent: 0, SendRate: 0}
            return
          } else {
            //increment timer
          }
        } else {
          //increment timer
        }
    }
  }
}
