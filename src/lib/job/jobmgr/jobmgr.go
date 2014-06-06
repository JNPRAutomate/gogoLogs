package jobmgr

import (
  "lib/job/jobmsg"
  "lib/job"
  "log"
)

type JobMgr struct {
  MainChannel chan job.Job
  CtrlChannel chan jobmsg.JobMsg
  JobHooks map[string]JobHook
}

type JobHook struct {
  ID string
  CtrlChannel chan jobmsg.JobMsg
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
        }
        jm.JobHooks[newJob.ID]  = newJobHook
        go newJob.Start()
      case newJobMsg := <- jm.CtrlChannel:
        log.Println("IIIIIIIIIIIIIIIIII",newJobMsg.ID)
        log.Println("HHHHHHHHHHHHHh",jm.JobHooks[newJobMsg.ID])
        jm.JobHooks[newJobMsg.ID].CtrlChannel <- jobmsg.JobMsg{Action:jobmsg.Stop}
        //remove job hook from slice
        delete(jm.JobHooks,newJobMsg.ID)
      default:
        continue
    }
  }
}

func(jm *JobMgr) removeJobHook(id *string) error {
  //remove job from map
  return nil
}
