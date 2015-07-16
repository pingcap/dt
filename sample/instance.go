package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/ngaut/log"
)

var (
	logDir = flag.String("log-dir", "/tmp/dt/instance", "the directory to store log")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	os.MkdirAll(*logDir, 0755)
	log.SetLevelByString("debug")
	log.SetOutputByName(*logDir + "/inst.log")

	go runHTTPProbeResult()

	l, err := net.Listen("tcp", ":54300")
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err = l.Accept()
		if err != nil {
			log.Warning("accept err - ", err)
		}
		// TODO: do sth
	}
}

func runHTTPProbeResult() {
	m := mux.NewRouter()
	m.HandleFunc("/probe/server/start", probeTest)
	m.HandleFunc("/probe/server/restart", probeTest)
	m.HandleFunc("/probe/server/dropport", probeTest)
	m.HandleFunc("/probe/server/recoverport", probeTest)
	m.HandleFunc("/probe/server/pause", probeTest)
	m.HandleFunc("/probe/server/continue", probeTest)
	m.HandleFunc("/probe/server/stop", probeTest)

	http.Handle("/", m)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("listen failed, err:", err)
	}
}

func probeTest(w http.ResponseWriter, r *http.Request) {
	log.Info("probe")
	w.WriteHeader(http.StatusOK)
}
