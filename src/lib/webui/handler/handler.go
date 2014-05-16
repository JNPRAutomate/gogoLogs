package handlers

import (
	"fmt"
	"net/http"
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
func (h *Handler) startJob(){
	//start new log sending task
	//issue via the jobs channel
	//return job was success or failure
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
	reqRouter.HandleFunc("/job/start/{ID}",h.startJob)
	reqRouter.HandleFunc("/job/stop/{ID}",h.startJob)
	reqRouter.HandleFunc("/job/status/{ID}",h.startJob)
	reqRouter.HandleFunc("/job/stats/{ID}",h.startJob)
	reqRouter.HandleFunc("/",h.manage)
	http.Handle("/",reqRouter)
	http.ListenAndServe(":"+strconv.Itoa(h.HttpPort),nil)
}
