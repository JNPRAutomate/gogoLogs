package handlers

import (
	"fmt"
	"net/http"
	"lib/job"
	"github.com/gorilla/mux"
)

/*Handler - Creates and managers WebUI */
type Handler struct {
	HttpPort int
	LogDir string
	jobChan chan
	statsChan chan
}

/*NewHandler creates new handler and returns in */
func NewHandler(jc chan, sc chan, p int, ld string) Handler {
	h := Handler{
		HttpPort:p,
		jobChan:jc,
		statsChan:sc,
		LogDir: ld
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

	//int values
	var rate string
	var jobRate int

	var syslogFacility string
	var jobSyslogFacility int

	var syslogPriority string
	var jobSyslogPriority int

	//string values
	var destHost string
	var logFileName string
	var protocol string
	var sourceHost string



	if sourceHost = req.FormValue("sourceHost"); sourceHost == nil {

	}
	if syslogFacility = req.FormValue("syslogFacility"); sourceHost == nil {
		if jobSyslogFacility , err = strconv.Atoi; err != nil {

		}
	}
	if syslogPriority = req.FormValue("syslogPriority"); sourceHost == nil {
		if jobSyslogPriority , err = strconv.Atoi; err != nil {

		}
	}
	if destHost = req.FormValue("destHost"); sourceHost == nil {

	}
	if logFileName = req.FormValue("logFileName"); sourceHost == nil {

	}
	if protocol = req.FormValue("protocol"); sourceHost == nil {

	}
	if rate = req.FormValue("rate"); rate == nil {
		if jobRate, err = strconv.Atoi(rate); err != nil {

		}
	}
	var job job.Job
	job,err = job.NewJob(jobRate,jobSyslogFacility,jobSyslogPriority,logFileName,h.jobChan)
}

/*stopJob stops log job*/
func (h *Handler) stopJob(){
	//request job to stop based on ID
	//return if stop was success or failure
}

/*statusJob checks job status*/
func (h *Handler) statusJob(){
	//check on job with ID
	//return job status
}

/*statsJob pulls stats for current job*/
func (h *Handler) statsJob(){
	//check for stats on job
	//return stats struct
}

/*manage Main UI access */
func (h *Handler) manage() {
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
	reqRouter.HandleFunc("/job/stop/{ID}",h.startJob)
	reqRouter.HandleFunc("/job/status/{ID}",h.startJob)
	reqRouter.HandleFunc("/job/stats/{ID}",h.startJob)
	reqRouter.HandleFunc("/",h.manage)
	http.Handle("/",reqRouter)
	http.ListenAndServe(":"+strconv.Itoa(h.HttpPort),nil)
}
