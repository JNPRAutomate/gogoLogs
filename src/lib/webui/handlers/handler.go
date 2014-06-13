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
	"encoding/json"
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
	name = h.logFiles[id].Path;
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
	IDint, _ := strconv.Atoi(ID)
	//send message to stop job
	newMsg := jobmsg.JobMsg{ID:ID,Action:jobmsg.Stop}
	h.ctrlChan <- newMsg
	fmt.Fprintf(w, "{\"id\":%d}", IDint)
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
	vars := mux.Vars(req)
	ID := vars["ID"]
	IDint, _ := strconv.Atoi(ID)
	fmt.Fprintf(w, "{\"id\":%d,\"count\":%d,\"rate\":%d}", IDint,h.jobMgr.Stats[ID].Count,h.jobMgr.Stats[ID].Rate)
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
		<style>
			.tdalign {
				text-align: left;
				width: 50px;
			}
			input {
				width: 75px
			}
			table {
				padding: 20px;
			}
			td {
				left-padding: 10px;
				right-padding: 10px;
			}
			.logRow {
				background-color: #cccccc;
			}
			.title {
				text-align: center;
			}
			.red {
				background-color: red;
			}
			.green {
				background-color: green;
			}
		</style>
	</head>
	<body>
		<h1 class="title">gogoLogs - A log sending tool</h1>
		<table>
			<tr><th>File Name</th><th>Status</th><th>Settings</th><th>Stats</th></tr>
			{{ range . }}
				<tr class="logRow" id="item{{ .ID }}">
					<td class="tdalign">{{ .Info.Name }} </td>
					<td id="status{{ .ID }}" class="red"></td>
					<td>
						<label for="rate{{.ID}}">Rate</label> <input type="text" id="rate{{ .ID }}" >
						<label for="syslogFacility{{.ID}}">Syslog Facility</label> <input type="text" id="syslogFacility{{.ID}}">
						<label for="syslogPriority{{.ID}}">Syslog Priority</label> <input type="text" id="syslogPriority{{.ID}}">
						<label for="destHost{{.ID}}">Dest IP</label> <input type="text" id="destHost{{.ID}}">
						<label for="sourceHost{{.ID}}">Source Host Name</label> <input type="text" id="sourceHost{{.ID}}">
						<button id="item{{.ID}}" type="button" data-type="start" data-id="{{.ID}}">Start</button>
						<button id="item{{.ID}}" type="button" data-type="stop" data-id="{{.ID}}">Stop</button>
					</td>
					<td id="stats{{.ID}}">None</td>
				</tr>
			{{ end }}

			<script src="//localhost/jquery.js"></script>
			<script type="application/javascript">
				window.onload = function() {

					var timeoutID = window.setInterval(function(){
						$.ajax({
							type: "GET",
							url: "job/list",
							dataType:"json",
							success: function(data,textStatus,jqxhr) {
								for (var i = 0; i < data.jobList.length; i++) {
									$("#status" + data.jobList[i]).removeClass("red");
									$("#status" + data.jobList[i]).addClass("green");
									$.ajax({
										type:"GET",
										url:"job/stats/" + data.jobList[i],
										dataType: "json",
										success: function(data,textStatus,jqxhr) {
											$("#stats" + data.id).empty().append("<div>Count: " + data.count + " Rate: " + data.rate + " </div>");
										}
									})
								}
							}
						});
					},3000);

					$( "button[id*='item']" ).click(function() {
						var itemid = $(this).data("id");
						if ($(this).data("type") === "start") {
							$("#status" + itemid).removeClass("red");
							$("#status" + itemid).addClass("green");
							$.ajax({
								type: "POST",
								url: "job/start",
								data: {
									rate:$("#rate"+itemid).val() ,
									syslogFacility: $("#syslogFacility"+itemid).val(),
									syslogPriority: $("#syslogPriority"+itemid).val(),
									destHost: $("#destHost"+itemid).val(),
									logFileID: $(this).data("id"),
									protocol: "udp",
									sourceHost: $("#sourceHost"+itemid).val()
								},
								success: function(data,textStatus,jqxhr){},
								dataType: "json"
							});
						} else if ($(this).data("type") === "stop") {
							var itemid = $(this).data("id");
							$.ajax({
								type: "POST",
								url: "job/stop/" + itemid,
								success: function(data,textStatus,jqxhr){
									$("#status" + data.id).removeClass("green");
									$("#status" + data.id).addClass("red");
									$("#stats" + data.id).empty().append("<div>None</div>");
								},
								dataType: "json"
							});
						}
						return false;
					});
					//Loop and get current stats
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

func (h *Handler) jobList(w http.ResponseWriter, req *http.Request) {
	//grab the current list of job IDs as an array
  jobList := make([]string,0)
	for key,_ := range h.jobMgr.JobHooks {
		jobList = append(jobList,key)
	}
	jsonData, _ := json.Marshal(jobList)
	fmt.Fprintf(w, "{\"jobList\":%s}", jsonData)
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
	reqRouter.HandleFunc("/job/list",h.jobList)
	reqRouter.HandleFunc("/",h.manage)
	http.Handle("/",reqRouter)
	go h.jobMgr.Run()
	http.ListenAndServe(":"+strconv.Itoa(h.HttpPort),nil)
}
