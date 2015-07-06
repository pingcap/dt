package agent

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

func runHTTPServer(a *Agent) error {
	log.Debug("start: runHTTPServer")
	m := mux.NewRouter()

	m.HandleFunc("/api/instance/start", a.inst.apiStart).Methods("Post", "Put")
	m.HandleFunc("/api/instance/restart", a.inst.apiRestart).Methods("Post", "Put")
	m.HandleFunc("/api/instance/pause", a.inst.apiPause).Methods("Post", "Put")
	m.HandleFunc("/api/instance/continue", a.inst.apiContinue).Methods("Post", "Put")
	m.HandleFunc("/api/instance/stop", a.inst.apiStop).Methods("Post", "Put")
	m.HandleFunc("/api/instance/dropport", a.inst.apiDropPort).Methods("Post", "Put")
	m.HandleFunc("/api/instance/recoverport", a.inst.apiRecoverPort).Methods("Post", "Put")
	m.HandleFunc("/api/instance/backupdata", a.inst.apiBackupData).Methods("Post", "Put")
	m.HandleFunc("/api/instance/cleanupdata", a.inst.apiCleanUpData).Methods("Post", "Put")
	m.HandleFunc("/api/agent/shutdown", a.apiShutdown).Methods("Post", "Put")

	http.Handle("/", m)
	err := http.ListenAndServe(a.Addr, nil)
	if err != nil {
		a.exitCh <- err
	}

	return errors.Trace(err)

}

func (inst *Instance) apiStart(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: apiStartInstance")
	cmd := r.FormValue("cmd")
	// probe := r.FormValue("probe")
	inst.dataDir = r.FormValue("dir")

	if err := inst.Start(cmd); err != nil {
		util.WriteHTTPError(w, "start instance failed, err:"+err.Error())
		return
	}
	// TODO: add probe

	w.WriteHeader(http.StatusOK)
}

// TODO: implement
func (inst *Instance) apiRestart(w http.ResponseWriter, r *http.Request)     {}
func (inst *Instance) apiStop(w http.ResponseWriter, r *http.Request)        {}
func (inst *Instance) apiPause(w http.ResponseWriter, r *http.Request)       {}
func (inst *Instance) apiContinue(w http.ResponseWriter, r *http.Request)    {}
func (inst *Instance) apiBackupData(w http.ResponseWriter, r *http.Request)  {}
func (inst *Instance) apiCleanUpData(w http.ResponseWriter, r *http.Request) {}
func (inst *Instance) apiDropPort(w http.ResponseWriter, r *http.Request)    {}
func (inst *Instance) apiRecoverPort(w http.ResponseWriter, r *http.Request) {}

func (a *Agent) apiShutdown(w http.ResponseWriter, r *http.Request) {
	if err := a.Shutdown(); err != nil {
		util.WriteHTTPError(w, "shutdown agent failed, err:"+err.Error())
	}

	w.WriteHeader(http.StatusOK)
}
