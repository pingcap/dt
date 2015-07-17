package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

var (
	addr = flag.String("addr", ":9090", "http listen addr")
)

func main() {
	log.Debug("start probe")
	runHTTPProbeResult()
}

func runHTTPProbeResult() {
	m := mux.NewRouter()
	m.HandleFunc("/probe/server/start", probeStart)
	m.HandleFunc("/probe/server/init", probeTest)
	m.HandleFunc("/probe/server/restart", probeTest)
	m.HandleFunc("/probe/server/dropport", probeTest)
	m.HandleFunc("/probe/server/recoverport", probeTest)
	m.HandleFunc("/probe/server/pause", probeTest)
	m.HandleFunc("/probe/server/continue", probeTest)
	m.HandleFunc("/probe/server/stop", probeTest)

	http.Handle("/", m)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("listen failed, err:", err)
	}
}

func probeStart(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probestart")
	// TODO : do more detailed tests
	arg := "./cockroach kv put a 1"
	if _, err := util.ExecCmd(arg, os.Stdout); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Info("end: probestart")
}

func probeTest(w http.ResponseWriter, r *http.Request) {
	log.Info("probe")
	w.WriteHeader(http.StatusOK)
}
