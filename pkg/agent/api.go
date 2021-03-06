package agent

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

func runHTTPServer(a *Agent) error {
	log.Debug("start: runHTTPServer")
	m := mux.NewRouter()
	inst := a.inst

	m.HandleFunc("/api/instance/start", inst.apiStart).Methods("Post", "Put")
	m.HandleFunc("/api/instance/set", inst.apiSet).Methods("Post", "Put")
	m.HandleFunc("/api/instance/restart", inst.apiRestart).Methods("Post", "Put")
	m.HandleFunc("/api/instance/pause", inst.apiPause).Methods("Post", "Put")
	m.HandleFunc("/api/instance/continue", inst.apiContinue).Methods("Post", "Put")
	m.HandleFunc("/api/instance/stop", inst.apiStop).Methods("Post", "Put")
	m.HandleFunc("/api/instance/droppkg", inst.apiDropPkg).Methods("Post", "Put")
	m.HandleFunc("/api/instance/limitspeed", inst.apiLimitSeep).Methods("Post", "Put")
	m.HandleFunc("/api/instance/dropport", inst.apiDropPort).Methods("Post", "Put")
	m.HandleFunc("/api/instance/recoverport", inst.apiRecoverPort).Methods("Post", "Put")
	m.HandleFunc("/api/instance/backupdata", a.apiBackupData).Methods("Post", "Put")
	m.HandleFunc("/api/instance/cleanupdata", a.apiCleanupData).Methods("Post", "Put")
	m.HandleFunc("/api/agent/shutdown", a.apiShutdown).Methods("Post", "Put")

	http.Handle("/", m)
	err := http.ListenAndServe(a.Addr, nil)
	if err != nil {
		a.exitCh <- err
	}

	return errors.Trace(err)

}

func (inst *Instance) apiStart(w http.ResponseWriter, r *http.Request) {
	log.Info("api")
	cmd := r.FormValue("cmd")
	name := r.FormValue("name")
	if util.CheckIsEmpty(cmd, name) {
		util.RespHTTPErr(w, http.StatusBadRequest, "")
		return
	}

	if err := inst.Start(cmd, name); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("start instance failed, err - %v", err))
		return
	}
	time.Sleep(2 * time.Second)
	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiSet(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	if util.CheckIsEmpty(cmd) {
		util.RespHTTPErr(w, http.StatusBadRequest, "")
		return
	}

	if err := inst.Set(cmd); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("set instance failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiRestart(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	name := r.FormValue("name")
	if util.CheckIsEmpty(cmd, name) {
		util.RespHTTPErr(w, http.StatusBadRequest, "")
		return
	}

	if err := inst.Restart(cmd, name); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("restart instance, failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiPause(w http.ResponseWriter, r *http.Request) {
	if err := inst.Pause(); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("pause instance failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiContinue(w http.ResponseWriter, r *http.Request) {
	if err := inst.Continue(); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("continue instance failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiStop(w http.ResponseWriter, r *http.Request) {
	if err := inst.Stop(); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("stop instance failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiDropPkg(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	chain := r.FormValue("chain")
	percent, err := strconv.Atoi(r.FormValue("percent"))
	if util.CheckIsEmpty(port, chain) || err != nil {
		util.RespHTTPErr(w, http.StatusBadRequest,
			fmt.Sprintf("%v", err))
		return
	}

	if err := inst.DropPkg(chain, port, percent); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("drop instance port failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiLimitSeep(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	chain := r.FormValue("chain")
	unit := r.FormValue("unit")
	pkgs, err := strconv.Atoi(r.FormValue("pkgs"))
	if util.CheckIsEmpty(port, chain, unit) || err != nil {
		util.RespHTTPErr(w, http.StatusBadRequest,
			fmt.Sprintf("%v", err))
		return
	}

	if err := inst.LimitSpeed(chain, port, unit, pkgs); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("drop instance port failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiDropPort(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	if util.CheckIsEmpty(port) {
		util.RespHTTPErr(w, http.StatusBadRequest, "")
		return
	}

	if err := inst.DropPort(port); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("drop instance port failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (inst *Instance) apiRecoverPort(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	if util.CheckIsEmpty(port) {
		util.RespHTTPErr(w, http.StatusBadRequest, "")
		return
	}

	if err := inst.RecoverPort(port); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("recover instance port failed, err - %v", err))
		return
	}

	if err := ProbeResult(r.FormValue("probe")); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Agent) apiBackupData(w http.ResponseWriter, r *http.Request) {
	dstPath := r.FormValue("dir")
	if util.CheckIsEmpty(dstPath) {
		util.RespHTTPErr(w, http.StatusBadRequest, "")
		return
	}

	if err := a.BackupData(dstPath); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("backup instance data failed, err - %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Agent) apiCleanupData(w http.ResponseWriter, r *http.Request) {
	if err := a.CleanupData(); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("clean up instance data failed, err - %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Agent) apiShutdown(w http.ResponseWriter, r *http.Request) {
	if err := a.Shutdown(); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("shutdown agent failed, err - %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ProbeResult(url string) error {
	if url == "" {
		return nil
	}

	err := util.HTTPCall(url, "", nil)
	if err != nil {
		return err
	}

	return nil
}
