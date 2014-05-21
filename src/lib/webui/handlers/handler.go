package handlers

import (
	"fmt"
	"log"
	"net/http"
	"lib/job/jobmgr"
	"lib/job/jobmsg"
	"lib/job"
	"strconv"
	"lib/stats"
	"github.com/gorilla/mux"
)

/*Handler - Creates and managers WebUI */
type Handler struct {
	HttpPort int
	LogDir string
	jobChan chan job.Job
	ctrlChan chan jobmsg.JobMsg
	statsChan chan stats.Stats
	jobMgr jobmgr.JobMgr
}

/*NewHandler creates new handler and returns in */
func NewHandler(cc chan jobmsg.JobMsg, sc chan stats.Stats, p int, ld string) Handler {
	jc := make(chan job.Job,0)
	h := Handler{
		HttpPort:p,
		ctrlChan:cc,
		statsChan:sc,
		jobChan: jc,
		jobMgr: jobmgr.NewJobMgr(jc, cc),
		LogDir:ld,
	}
	return h
}

func (h *Handler) SetLogDir() {
	//change log dir
	//reinit the log struct
}

/*startJob starts new log job*/
func (h *Handler) startJob(w http.ResponseWriter, req *http.Request){
	//start new log sending task
	//issue via the jobs channel
	//return job was success or failure

	//create new job based upon request
	//start job as go routine
	var err error
	//int values
	var rate string
	var jobRate int

	var syslogFacility string
	var jobSyslogFacility int

	var syslogSeverity string
	var jobsyslogSeverity int

	//string values
	var destHost string
	var logFileName string
	var protocol string
	var sourceHost string
	var port string
	port = "514"


	if sourceHost = req.FormValue("sourceHost"); sourceHost != "" {

	}
	if syslogFacility = req.FormValue("syslogFacility"); syslogFacility != "" {
		if jobSyslogFacility , err = strconv.Atoi(syslogFacility); err != nil {

		}
	}
	if syslogSeverity = req.FormValue("syslogSeverity"); syslogSeverity != "" {
		if jobsyslogSeverity , err = strconv.Atoi(syslogSeverity); err != nil {

		}
	}
	if destHost = req.FormValue("destHost"); destHost != "" {

	}
	if logFileName = req.FormValue("logFileName"); logFileName != "" {

	}
	if protocol = req.FormValue("protocol"); protocol != "" {

	}
	if rate = req.FormValue("rate"); rate != "" {
		if jobRate, err = strconv.Atoi(rate); err != nil {

		}
	}
	var newJob job.Job
	if newJob,err = job.NewJob(&destHost,&port,&jobRate,&jobSyslogFacility,&jobsyslogSeverity,&sourceHost,&logFileName,h.ctrlChan); err != nil {
		log.Println(err,logFileName)
	} else {
		ID := newJob.GenID()
		fmt.Println(ID)
		h.jobMgr.MainChannel <-newJob
	}
	fmt.Fprintf(w, "Hi there, I love %s!", req.URL.Path[1:])
}

/*stopJob stops log job*/
func (h *Handler) stopJob(w http.ResponseWriter, req *http.Request){
	//request job to stop based on ID
	//return if stop was success or failure
	vars := mux.Vars(req)
	ID := vars["ID"]
	fmt.Println(ID)
	//send message to stop job
	h.jobMgr.CtrlChannel <- jobmsg.JobMsg{ID:ID,Action:jobmsg.Stop}
	fmt.Fprintf(w, "Hi there, I love %s!", req.URL.Path[1:])
}

/*statusJob checks job status*/
func (h *Handler) statusJob(w http.ResponseWriter, req *http.Request){
	//check on job with ID
	//return job status
}

/*statsJob pulls stats for current job*/
func (h *Handler) statsJob(w http.ResponseWriter, req *http.Request){
	//check for stats on job
	//return stats struct
}

/*manage Main UI access */
func (h *Handler) manage(w http.ResponseWriter, req *http.Request) {
	//handle access to basic UI
	//WebUI Mode
	//list files in specified directory

	//UI lines
	// DEST HOST, LOG List combo box, syslog facility (combo), syslog priority (combo), src host, protocol, rate, stop and start toggle button
}

/*Start starts the webUI handler*/
func (h *Handler) Start() {
	//setup handlers and listeners
	reqRouter := mux.NewRouter()
	reqRouter.HandleFunc("/job/start",h.startJob)
	reqRouter.HandleFunc("/job/stop/{ID}",h.stopJob)
	reqRouter.HandleFunc("/job/status/{ID}",h.statusJob)
	reqRouter.HandleFunc("/job/stats/{ID}",h.statsJob)
	reqRouter.HandleFunc("/",h.manage)
	http.Handle("/",reqRouter)
	go h.jobMgr.Run()
	http.ListenAndServe(":"+strconv.Itoa(h.HttpPort),nil)
}
