package jobmsg

const (
  start = 1
  stop = 2
  restart = 3
)

type JobMsg stuct {
  ID int //job that message is sent to
  action int //action to send to job
}
