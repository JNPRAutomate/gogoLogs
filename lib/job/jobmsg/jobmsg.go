package jobmsg

const (
  Start = 1
  Stop = 2
  Restart = 3
)

type JobMsg struct {
  ID string
  Action int //action to send to job
}
