package instance_agent

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/pingcap/dt/pkg/util"
)

func runHttpServer(a *Agent) error {
	m := mux.NewRouter()

	m.HandleFunc("/"+util.ActionStartInstance, a.apiStartInstance).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionRestartInstance, a.apiRestartInstance).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionPauseInstance, a.apiPauseInstance).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionContinueInstance, a.apiContinueInstance).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionStopInstance, a.apiStopInstance).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionDropPortInstance, a.apiDropPortInstance).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionRecoverPortInstance, a.apiRecoverPortInstance).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionBackupInstanceData, a.apiBackupInstanceData).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionCleanUpInstanceData, a.apiCleanUpInstanceData).Methods("Post", "Put")
	m.HandleFunc("/"+util.ActionShutdown, a.apiShutdown).Methods("Post", "Put")

	http.Handle("/", m)

	return http.ListenAndServe(a.Addr, nil)
}

func (a *Agent) apiStartInstance(w http.ResponseWriter, r *http.Request) {
	// TODO: args format
	args := r.FormValue("args")

	if err := a.StartInstance(args); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "start instance failed, err:"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

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
