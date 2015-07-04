package controller

import (
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/ngaut/log"
)

var ErrRegisterArgs = errors.New("invalid register args")

func runHttpServer(addr string, ctrl *Controller) {
	log.Debug("start: runHttpServer")
	m := mux.NewRouter()

	m.HandleFunc("/api/agent/register", ctrl.apiRegisterAgent).Methods("POST", "PUT")

	http.Handle("/", m)
	http.ListenAndServe(addr, nil)
}

func (ctrl *Controller) apiRegisterAgent(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: apiRegisterAgent")
	agentAddr := r.FormValue("addr")

	if agentAddr == "" {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, ErrRegisterArgs.Error())
		return
	}
	log.Info("apiRegisterAgent, info:", agentAddr)
	ctrl.agentInfoCh <- agentAddr

	w.WriteHeader(http.StatusOK)
}
