package agent

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/ngaut/log"
)

func runHttpServer(a *Agent) error {
	log.Debug("start: runHttpServer")
	m := mux.NewRouter()

	m.HandleFunc("/api/instance/start", a.apiStartInstance).Methods("Post", "Put")
	m.HandleFunc("/api/instance/restart", a.apiRestartInstance).Methods("Post", "Put")
	m.HandleFunc("/api/instance/pause", a.apiPauseInstance).Methods("Post", "Put")
	m.HandleFunc("/api/instance/continue", a.apiContinueInstance).Methods("Post", "Put")
	m.HandleFunc("/api/instance/stop", a.apiStopInstance).Methods("Post", "Put")
	m.HandleFunc("/api/instance/dropport", a.apiDropPortInstance).Methods("Post", "Put")
	m.HandleFunc("/api/instance/recoverport", a.apiRecoverPortInstance).Methods("Post", "Put")
	m.HandleFunc("/api/instance/backupdata", a.apiBackupInstanceData).Methods("Post", "Put")
	m.HandleFunc("/api/instance/cleanupdata", a.apiCleanUpInstanceData).Methods("Post", "Put")
	m.HandleFunc("/api/agent/shutdown", a.apiShutdown).Methods("Post", "Put")

	http.Handle("/", m)

	return http.ListenAndServe(a.Addr, nil)
}

func (a *Agent) apiStartInstance(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: apiStartInstance")
	// TODO: args format
	args := r.FormValue("args")

	if err := a.StartInstance(args); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "start instance failed, err:"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Info("end: apiStartInstance")

	return
}

// TODO: implement
func (a *Agent) apiRestartInstance(w http.ResponseWriter, r *http.Request)     {}
func (a *Agent) apiStopInstance(w http.ResponseWriter, r *http.Request)        {}
func (a *Agent) apiPauseInstance(w http.ResponseWriter, r *http.Request)       {}
func (a *Agent) apiContinueInstance(w http.ResponseWriter, r *http.Request)    {}
func (a *Agent) apiBackupInstanceData(w http.ResponseWriter, r *http.Request)  {}
func (a *Agent) apiCleanUpInstanceData(w http.ResponseWriter, r *http.Request) {}
func (a *Agent) apiShutdown(w http.ResponseWriter, r *http.Request)            {}
func (a *Agent) apiDropPortInstance(w http.ResponseWriter, r *http.Request)    {}
func (a *Agent) apiRecoverPortInstance(w http.ResponseWriter, r *http.Request) {}
