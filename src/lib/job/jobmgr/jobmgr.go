package jobmgr

import (
  "lib/job/jobmsg"
  "lib/job"
  "log"
)

type JobMgr struct {
  MainChannel chan job.Job
  CtrlChannel chan jobmsg.JobMsg
  JobHooks []JobHook
}

type JobHook struct {
  ID string
  CtrlChannel chan jobmsg.JobMsg
}

func NewJobMgr( mc chan job.Job, cc chan jobmsg.JobMsg) JobMgr {
  jm := JobMgr{MainChannel: mc, CtrlChannel: cc}
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
        jm.JobHooks = append(jm.JobHooks,newJobHook)
        go newJob.Start()
      case newJobMsg := <- jm.CtrlChannel:
        newJobHook := jm.findJob(&newJobMsg.ID) //change me to a map!
        newJobHook.CtrlChannel <- jobmsg.JobMsg{Action:jobmsg.Stop}
        //remove job hook from slice
      default:
        continue
    }
  }
}

func(jm *JobMgr) removeJobHook(id *string) error {
  //remove job from map
  return nil
}

func(jm *JobMgr) findJob(id *string) *JobHook {
  for item := range jm.JobHooks {
    log.Println("Item # ",jm.JobHooks[item].ID)
    if jm.JobHooks[item].ID == *id {
      log.Println("Item # ",jm.JobHooks[item].ID)
      return &jm.JobHooks[item]
    }
  }
  return &JobHook{ID:"",CtrlChannel:make(chan jobmsg.JobMsg,4096)}
}
