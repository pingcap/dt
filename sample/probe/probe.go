package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/cockroachdb/cockroach/client"
	"github.com/gorilla/mux"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

var (
	addr  = flag.String("addr", ":9090", "http listen addr")
	sAddr = flag.String("s-addr", "xia-pc:8080", "server addr")
)

var keyGlobal = 1

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Debug("start probe")
	runHTTPProbeResult()
}

func makeDBClient(addr string) *client.DB {
	var err error
	url := fmt.Sprintf("https://root@%s", addr)
	log.Info("url:", url)

	DB, err := client.Open(url)
	if err != nil {
		log.Fatal(err)
	}

	return DB
}

func runHTTPProbeResult() {
	m := mux.NewRouter()
	m.HandleFunc("/probe/server/start", probeStart)
	m.HandleFunc("/probe/server/init", probeTest)
	m.HandleFunc("/probe/server/restart", probeTest)
	m.HandleFunc("/probe/server/dropport", probeDropPort)
	m.HandleFunc("/probe/server/recoverport", probeRecoverPort)
	m.HandleFunc("/probe/server/pause", probePause)
	m.HandleFunc("/probe/server/continue", probeContinue)
	m.HandleFunc("/probe/server/stop", probeTest)

	http.Handle("/", m)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("listen failed, err:", err)
	}
}

func generateKey() string {
	keyGlobal++

	return fmt.Sprintf("%08d", keyGlobal)
}

func isInvertRet(result string) bool {
	if result == "pass" {
		return true
	}

	return false
}

func probePass(w http.ResponseWriter, isInvert, key, flag string) {
	exitCh := make(chan bool, 1)
	go func() {
		DB := makeDBClient(*sAddr)
		if err := DB.Put(key, flag+key); err != nil {
			util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
			exitCh <- true
		}
		exitCh <- false
	}()

	timeout := time.After(5 * time.Second)
	select {
	case exit := <-exitCh:
		if (exit && isInvertRet(isInvert)) || (!exit && !isInvertRet(isInvert)) {
			return
		}
		break
	case <-timeout:
		if !isInvertRet(isInvert) {
			break
		}
		util.RespHTTPErr(w, http.StatusInternalServerError, "timeout")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func probeStart(w http.ResponseWriter, r *http.Request) {
	key := generateKey()
	DB := makeDBClient(*sAddr)
	if err := DB.Put(key, "start"+key); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func probeDropPort(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe drop port")
	isInvert := r.FormValue("result")
	key := generateKey()

	probePass(w, isInvert, key, "dorpport")
	log.Debug("end: probe drop port")
}

func probeRecoverPort(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe recover port")
	isInvert := r.FormValue("result")
	key := generateKey()

	probePass(w, isInvert, key, "recoverport")
	log.Debug("end: probe recover port")
}

func probePause(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe pause")
	isInvert := r.FormValue("result")
	key := generateKey()

	probePass(w, isInvert, key, "pause")
	log.Debug("end: probe pause")
}

func probeContinue(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe continue")
	isInvert := r.FormValue("result")
	key := generateKey()

	probePass(w, isInvert, key, "continue")
	log.Debug("end: probe continue")
}

func probeTest(w http.ResponseWriter, r *http.Request) {
	log.Info("probe")
	w.WriteHeader(http.StatusOK)
}
