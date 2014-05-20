package jobmgr

import (
  "lib/job/jobmsg"
  "lib/job"
)

type JobMgr struct {
  MainChannel chan job.Job
  ctrlChannel chan jobmsg.JobMsg
  Jobs []JobHook
}

type JobHook struct {
  ID string
  ctrlChannel chan jobmsg.JobMsg
}

func NewJobMgr( mc chan job.Job, cc chan jobmsg.JobMsg) JobMgr {
  jm := JobMgr{MainChannel: mc, ctrlChannel: cc}
  return jm
}

func(jm *JobMgr) Run() {
  for {
    select {
      case newJob := <- jm.MainChannel:
        newJobID := newJob.GenID()
        newCtrlChan := make(chan jobmsg.JobMsg,0)
        newJobHook := JobHook{
          ID: newJobID,
          ctrlChannel:newCtrlChan,
        }
        jm.Jobs = append(jm.Jobs,newJobHook)
        go newJob.Start()
    }
  }
}

func(jm *JobMgr) Stop(id *string) error {
  stopJob := jm.findJob(id)
  stopJob.ctrlChannel <- jobmsg.JobMsg{Action:jobmsg.Stop}
  return nil
}

func(jm *JobMgr) removeJobHook(id *string) error {
  //remove job from map
  return nil
}

func(jm *JobMgr) findJob(id *string) JobHook {
  for item := range jm.Jobs {
    if &jm.Jobs[item].ID == id {
      return jm.Jobs[item]
    }
  }
  return JobHook{ID:""}
}
