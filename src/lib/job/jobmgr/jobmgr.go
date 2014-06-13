package jobmgr

import (
  "lib/job/jobmsg"
  "lib/stats/statsmsg"
  "lib/stats"
  "lib/job"
  "log"
)

type JobMgr struct {
  MainChannel chan job.Job
  CtrlChannel chan jobmsg.JobMsg
  StatsChannel chan statsmsg.StatsMsg
  JobHooks map[string] JobHook
  Stats map[string] stats.Stats
}

type JobHook struct {
  ID string
  CtrlChannel chan jobmsg.JobMsg
}

func NewJobMgr( mc chan job.Job, cc chan jobmsg.JobMsg) JobMgr {
  jm := JobMgr{MainChannel: mc, CtrlChannel: cc}
  jm.JobHooks = make(map[string]JobHook)
  jm.Stats = make(map[string]stats.Stats)
  jm.StatsChannel = make(chan statsmsg.StatsMsg,4096)
  return jm
}

func(jm *JobMgr) Run() {
  for {
    select {
      case newJob := <- jm.MainChannel:
        newJobHook := JobHook{
          ID: newJob.ID,
          CtrlChannel:newJob.CtrlChannel,
        }
        jm.JobHooks[newJob.ID] = newJobHook
        //create new stats history
        jm.Stats[newJob.ID] = stats.Stats{Count:0,Rate:0}
        newJob.StatsChannel = jm.StatsChannel
        go newJob.Start()
      case newJobMsg := <- jm.CtrlChannel:
        if (newJobMsg.Action == jobmsg.Stop) {
          if jm.JobHooks[newJobMsg.ID].ID != "" {
            jm.JobHooks[newJobMsg.ID].CtrlChannel <- jobmsg.JobMsg{Action:jobmsg.Stop}
          }
        }
        //remove job hook from slice
        delete(jm.JobHooks,newJobMsg.ID)
      case newStatsMsg := <- jm.StatsChannel:
        delete(jm.Stats,newStatsMsg.ID)
        jm.Stats[newStatsMsg.ID] = stats.Stats{Count:newStatsMsg.TotalSent,Rate:newStatsMsg.SendRate}
        log.Println(jm.Stats[newStatsMsg.ID])
    }
  }
}
