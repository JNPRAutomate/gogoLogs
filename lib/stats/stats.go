package stats

import (
  "time"
)

type Stats struct {
  Count uint
  Rate uint
}

type StatsHistory struct {
  ID string //ID of the Job with stats
  LogsSent []SendHistory //count of messages sent
  SendRate []RateHistory
}

type SendHistory struct {
  Time time.Time
  Count uint
}

type RateHistory struct {
  Time time.Time
  Rate uint
}

func NewSendHistory(c uint) SendHistory {
  sh := SendHistory{Time:time.Now(),Count:c}
  return sh
}

func NewRateHistory(r uint) RateHistory {
  rh := RateHistory{Time:time.Now(),Rate:r}
  return rh
}
