package jobmgr

import (
  "lib/job/jobmsg"
  "lib/stats/statsmsg"
  "lib/job"
)

type JobMgr struct {
  MainChannel chan job.Job
  CtrlChannel chan jobmsg.JobMsg
  JobHooks map[string]JobHook
}

type JobHook struct {
  ID string
  CtrlChannel chan jobmsg.JobMsg
  StatsChannel chan statsmsg.StatsMsg
}

func NewJobMgr( mc chan job.Job, cc chan jobmsg.JobMsg) JobMgr {
  jm := JobMgr{MainChannel: mc, CtrlChannel: cc}
  jm.JobHooks = make(map[string]JobHook)
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
        jm.JobHooks[newJob.ID]  = newJobHook
        go newJob.Start()
      case newJobMsg := <- jm.CtrlChannel:
        if (newJobMsg.Action == jobmsg.Stop) {
          jm.JobHooks[newJobMsg.ID].CtrlChannel <- jobmsg.JobMsg{Action:jobmsg.Stop}
        }
        //remove job hook from slice
        delete(jm.JobHooks,newJobMsg.ID)
      default:
        //collect all stats here on every tick
        continue
    }
  }
}

func(jm *JobMgr) removeJobHook(id *string) error {
  //remove job from map
  return nil
}
