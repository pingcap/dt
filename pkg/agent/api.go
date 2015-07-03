package instance_agent

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/pingcap/dt/pkg/util"
)

func runHttpServer(a *Agent) error {
	m := mux.NewRouter()

	m.HandleFunc(util.UrlStartInstance, a.apiStartInstance).Methods("Post", "Put")
	m.HandleFunc(util.UrlRestartInstance, a.apiRestartInstance).Methods("Post", "Put")
	m.HandleFunc(util.UrlPauseInstance, a.apiPauseInstance).Methods("Post", "Put")
	m.HandleFunc(util.UrlContinuePauseInstance, a.apiContinuePauseInstance).Methods("Post", "Put")
	m.HandleFunc(util.UrlStopInstance, a.apiStopInstance).Methods("Post", "Put")
	m.HandleFunc(util.UrlDropPortInstance, a.apiDropPortInstance).Methods("Post", "Put")
	m.HandleFunc(util.UrlContinueInstance, a.apiContinuePortInstance).Methods("Post", "Put")
	m.HandleFunc(util.UrlBackupInstanceData, a.apiBackupInstanceData).Methods("Post", "Put")
	m.HandleFunc(util.UrlCleanUpInstanceData, a.apiCleanUpInstanceData).Methods("Post", "Put")
	m.HandleFunc(util.UrlShutdown, a.apiShutdown).Methods("Post", "Put")

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
func (a *Agent) apiRestartInstance(w http.ResponseWriter, r *http.Request)       {}
func (a *Agent) apiStopInstance(w http.ResponseWriter, r *http.Request)          {}
func (a *Agent) apiPauseInstance(w http.ResponseWriter, r *http.Request)         {}
func (a *Agent) apiContinuePauseInstance(w http.ResponseWriter, r *http.Request) {}
func (a *Agent) apiBackupInstanceData(w http.ResponseWriter, r *http.Request)    {}
func (a *Agent) apiCleanUpInstanceData(w http.ResponseWriter, r *http.Request)   {}
func (a *Agent) apiShutdown(w http.ResponseWriter, r *http.Request)              {}
func (a *Agent) apiDropPortInstance(w http.ResponseWriter, r *http.Request)      {}
func (a *Agent) apiContinuePortInstance(w http.ResponseWriter, r *http.Request)  {}
