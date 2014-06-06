package handlers

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"lib/job/jobmgr"
	"lib/job/jobmsg"
	"lib/job"
	"strconv"
	"lib/stats"
	"html/template"
	"github.com/gorilla/mux"
	"path/filepath"
)

type File struct {
	Info os.FileInfo
	Path string
	ID int
}

/*Handler - Creates and managers WebUI */
type Handler struct {
	HttpPort int
	LogDir string
	logFiles []File
	jobChan chan job.Job
	ctrlChan chan jobmsg.JobMsg
	statsChan chan stats.Stats
	jobMgr jobmgr.JobMgr
}

/*NewHandler creates new handler and returns in */
func NewHandler(cc chan jobmsg.JobMsg, sc chan stats.Stats, p int, ld string) Handler {
	jc := make(chan job.Job,4096)
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

func (h *Handler) listFiles() {
	filepath.Walk(h.LogDir, func(path string, info os.FileInfo, err error) error {
    if (!info.IsDir()) {
			newFile := File{Info:info,Path:path,ID:len(h.logFiles)}
      h.logFiles = append(h.logFiles,newFile)
    }
    return nil
  })
}

func (h *Handler) logFileNameByID(id int) string {
	var name string
	name = h.logFiles[0].Path;
	return name
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
	var jobSyslogSeverity int

	//string values
	var destHost string
	var logFileID string
	var jobLogFileID int
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
		if jobSyslogSeverity , err = strconv.Atoi(syslogSeverity); err != nil {

		}
	}
	if destHost = req.FormValue("destHost"); destHost != "" {

	}
	if logFileID = req.FormValue("logFileID"); logFileID != "" {
		if jobLogFileID , err = strconv.Atoi(logFileID); err != nil {
			log.Println("EEEEEEEEEEEEEee",err)
		}
	}
	if protocol = req.FormValue("protocol"); protocol != "" {

	}
	if rate = req.FormValue("rate"); rate != "" {
		if jobRate, err = strconv.Atoi(rate); err != nil {

		}
	}
	var newJob job.Job
	//find the logfile name by ID
	logFileName := h.logFileNameByID(jobLogFileID)
	if newJob,err = job.NewJob(&destHost,&port,&jobRate,&jobSyslogFacility,&jobSyslogSeverity,&sourceHost,&logFileName,make(chan jobmsg.JobMsg,4096)); err != nil {
		log.Println(err,logFileName)
	} else {
		log.Println("SEEEEEEEEEEEEEEEEEEEEEEEEET",jobLogFileID)
		newJob.SetID(jobLogFileID)
		h.jobMgr.MainChannel <- newJob
	}
	fmt.Fprintf(w, "Hi there, I love %s!", req.URL.Path[1:])
}

/*stopJob stops log job*/
func (h *Handler) stopJob(w http.ResponseWriter, req *http.Request){
	//request job to stop based on ID
	//return if stop was success or failure
	vars := mux.Vars(req)
	ID := vars["ID"]
	//send message to stop job
	newMsg := jobmsg.JobMsg{ID:ID,Action:jobmsg.Stop}
	fmt.Println("MMMMMMMMMMMMMMMMMMM",newMsg)
	h.ctrlChan <- newMsg
	fmt.Fprintf(w, "Hi there, I love %s!", ID)
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
	var HTML_DATA string
	HTML_DATA = `
<!doctype html>
<html>
	<head>
		<title>Log Control</title>
		<meta charset="UTF-8">
	</head>
	<body>
		<table>
			<tr><th>File Name</th><th>Control</th><th>Stats</th></tr>
			{{ range . }}
			<tr id="item{{ .ID }}"><td style="text-align: left;">{{ .Info.Name }} </td><td><button id="item{{ .ID }}" type="button" data-type="start" data-id="{{.ID}}">Start</button> <button id="item{{ .ID }}" type="button" data-type="stop" data-id="{{.ID}}">Stop</button> <label for="rate{{.ID}}">Rate</label> <input type="text" name="rate{{ .ID }}" > <label for="syslogFacility{{.ID}}">Syslog Facility</label> <input type="text" name="syslogFacility{{.ID}}">  <label for="syslogPriority{{.ID}}">Syslog Priority</label> <input type="text" name="syslogPriority{{.ID}}"> <label for="destHost{{.ID}}">Dest IP</label> <input type="text" name="destHost{{.ID}}"> <label for="sourceHost{{.ID}}">Source Host Name</label> <input type="text" name="sourceHost{{.ID}}"> </td><td>Sent: 1234 Rate: 433/s</td></tr>
			{{ end}}
			<script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"></script>
			<script type="application/javascript">
				window.onload = function() {
					$( "button" ).click(function() {
						console.log($(this).data("type"));
						if ($(this).data("type") === "start") {
							$.ajax({
								type: "POST",
								url: "job/start",
								data: {rate:1 , syslogFacility:1 , syslogPriority:1 , destHost: "10.0.1.100", logFileID: $(this).data("id"), protocol: "udp", sourceHost: "appple"},
								success: function(data,textStatus,jqxhr){},
								dataType: "json"
							});
						} else if ($(this).data("type") === "stop") {
							$.ajax({
								type: "POST",
								url: "job/stop/" + $(this).data("id"),
								data: {rate:1 , syslogFacility:1 , syslogPriority:1 , destHost: "10.0.1.100", logFileID: $(this).data("id"), protocol: "udp", sourceHost: "appple"},
								success: function(data,textStatus,jqxhr){},
								dataType: "json"
							});
						}
						return false;
					});
				};
			</script>
	</body>
</html>

`
	//UI lines
	// DEST HOST, LOG List combo box, syslog facility (combo), syslog priority (combo), src host, protocol, rate, stop and start toggle button
	t := template.New("TEST TEMP")
	t, _ = t.Parse(HTML_DATA)
	t.Execute(w,h.logFiles)
}

/*Start starts the webUI handler*/
func (h *Handler) Start() {
	h.listFiles()
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
