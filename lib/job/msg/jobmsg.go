package jobmsg

const (
  start = 1
  stop = 2
  restart = 3
)

type JobMsg struct {
  ID int
  action int
}
