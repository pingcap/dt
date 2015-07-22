package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
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

const (
	timeoutFlag = true
	passResult  = "pass"
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

func checkResult(err error, result string, flag bool) bool {
	if flag == timeoutFlag {
		if result == "pass" {
			return false
		}
		return true
	}

	if (err == nil && result == passResult) || (err != nil && result != passResult) {
		return true
	}

	return false
}

func probePass(w http.ResponseWriter, isPass, key, flag string) {
	resultCh := make(chan error, 1)
	go func() {
		DB := makeDBClient(*sAddr)
		if err := DB.Put(key, flag+key); err != nil {
			resultCh <- err
		}
		resultCh <- nil
	}()

	timeout := time.After(5 * time.Second)
	select {
	case ret := <-resultCh:
		if checkResult(ret, isPass, !timeoutFlag) {
			break
		}
		if ret == nil {
			util.RespHTTPErr(w, http.StatusBadRequest, "")
			return
		}
		util.RespHTTPErr(w, http.StatusInternalServerError, ret.Error())
		return
	case <-timeout:
		if checkResult(nil, isPass, timeoutFlag) {
			break
		}
		util.RespHTTPErr(w, http.StatusInternalServerError, "timeout")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func probeStart(w http.ResponseWriter, r *http.Request) {
	key := generateKey()
	time.Sleep(5 * time.Second)

	DB := makeDBClient(*sAddr)
	if err := DB.Put(key, "start"+key); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func probeDropPort(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe drop port")
	ret := r.FormValue("result")
	timeout := r.FormValue("timeout")
	t, err := strconv.Atoi(timeout)
	if err != nil {
		util.RespHTTPErr(w, http.StatusBadRequest, err.Error())
	}
	key := generateKey()

	time.Sleep(time.Duration(t) * time.Second)
	probePass(w, ret, key, "dorpport")
	log.Debug("end: probe drop port")
}

func probeRecoverPort(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe recover port")
	ret := r.FormValue("result")
	if ret == "" {
		ret = passResult
	}
	key := generateKey()

	probePass(w, ret, key, "recoverport")
	log.Debug("end: probe recover port")
}

func probePause(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe pause")
	ret := r.FormValue("result")
	if ret == "" {
		ret = passResult
	}
	key := generateKey()

	probePass(w, ret, key, "pause")
	log.Debug("end: probe pause")
}

func probeContinue(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe continue")
	ret := r.FormValue("result")
	if ret == "" {
		ret = passResult
	}
	key := generateKey()

	probePass(w, ret, key, "continue")
	log.Debug("end: probe continue")
}

func probeTest(w http.ResponseWriter, r *http.Request) {
	log.Info("probe")
	w.WriteHeader(http.StatusOK)
}
