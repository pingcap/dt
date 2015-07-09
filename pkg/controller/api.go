package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

var ErrRegisterArgs = errors.New("invalid register args")

func runHTTPServer(addr string, ctrl *Controller) {
	log.Debug("start: runHTTPServer")
	m := mux.NewRouter()

	m.HandleFunc("/api/agent/register", ctrl.apiRegisterAgent).Methods("POST", "PUT")

	http.Handle("/", m)
	http.ListenAndServe(addr, nil)
}

func (ctrl *Controller) apiRegisterAgent(w http.ResponseWriter, r *http.Request) {
	agentAddr := r.FormValue("addr")

	if agentAddr == "" {
		util.RespHTTPErr(w, http.StatusInternalServerError,
			fmt.Sprintf("shutdown agent failed, err:%v", ErrRegisterArgs))
		return
	}
	log.Info("apiRegisterAgent, info:", agentAddr)
	ctrl.agentInfoCh <- agentAddr

	w.WriteHeader(http.StatusOK)
}
