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
  JobHooks map[string] JobHook
  Stats map[string] stats.Stats
}

type JobHook struct {
  ID string
  CtrlChannel chan jobmsg.JobMsg
  StatsChannel chan statsmsg.StatsMsg
}

func NewJobMgr( mc chan job.Job, cc chan jobmsg.JobMsg) JobMgr {
  jm := JobMgr{MainChannel: mc, CtrlChannel: cc}
  jm.JobHooks = make(map[string]JobHook)
  jm.Stats = make(map[string]stats.Stats)
  return jm
}

func(jm *JobMgr) Run() {
  for {
    select {
      case newJob := <- jm.MainChannel:
        newJobHook := JobHook{
          ID: newJob.ID,
          CtrlChannel:newJob.CtrlChannel,
          StatsChannel: newJob.StatsChannel,
        }
        jm.JobHooks[newJob.ID] = newJobHook
        //create new stats history
        jm.Stats[newJob.ID] = stats.Stats{Count:0,Rate:0}
        go newJob.Start()
      case newJobMsg := <- jm.CtrlChannel:
        if (newJobMsg.Action == jobmsg.Stop) {
          jm.JobHooks[newJobMsg.ID].CtrlChannel <- jobmsg.JobMsg{Action:jobmsg.Stop}
        }
        //remove job hook from slice
        delete(jm.JobHooks,newJobMsg.ID)
      default:
        for key, _ := range jm.JobHooks {
          select {
            case statsMsg := <- jm.JobHooks[key].StatsChannel:
              delete(jm.Stats,key)
              jm.Stats[key] = stats.Stats{Count:statsMsg.TotalSent,Rate:statsMsg.SendRate}
              log.Println(jm.Stats[key])
             default:
              continue
          }
        }
        //collect all stats here on every tick
        continue
    }
  }
}

func(jm *JobMgr) removeJobHook(id *string) error {
  //remove job from map
  return nil
}
