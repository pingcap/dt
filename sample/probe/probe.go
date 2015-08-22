package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/cockroach/client"
	"github.com/gorilla/mux"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

var (
	addr   = flag.String("addr", ":9090", "http listen addr")
	sAddrs = flag.String("s-addr", "xia-pc:8080,xia-pc:8081,xia-pc:8082", "server addrs")
)

const (
	passResult = "pass"
)

var (
	keyGlobal int64 = 1
	emptyKV         = client.KeyValue{}
	servAddrs []string
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Debug("start probe")
	servAddrs = strings.Split(*sAddrs, ",")
	if len(servAddrs) != 3 {
		log.Fatal("bad argments")
	}
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
	m.HandleFunc("/probe/server/restart", probeRestart)
	m.HandleFunc("/probe/server/dropport", probeDrop)
	m.HandleFunc("/probe/server/recoverport", probeRecoverPort)
	m.HandleFunc("/probe/server/pause", probePause)
	m.HandleFunc("/probe/server/continue", probeContinue)
	m.HandleFunc("/probe/server/stop", probeDrop)

	http.Handle("/", m)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("listen failed, err:", err)
	}
}

func getCurrentKey() string {
	return fmt.Sprintf("%08d", atomic.LoadInt64(&keyGlobal))
}

func generateKey() string {
	return fmt.Sprintf("%08d", atomic.AddInt64(&keyGlobal, 1))
}

func checkResult(err error, result string) error {
	if (err == nil && result == passResult) || (err != nil && result != passResult) {
		return nil
	}
	if err == nil && result != passResult {
		return errors.New("bad request")
	}

	return err
}

func probePass(isPass string, doOp func() (client.KeyValue, error)) (*client.KeyValue, error) {
	var kv client.KeyValue
	var err error
	resultCh := make(chan error, 1)

	go func() {
		if kv, err = doOp(); err != nil {
			resultCh <- err
		}
		resultCh <- nil
	}()

	timeout := time.After(5 * time.Second)
	select {
	case ret := <-resultCh:
		return &kv, checkResult(ret, isPass)
	case <-timeout:
		log.Warning("timeout:", isPass)
		return &kv, checkResult(errors.New("timeout"), isPass)
	}

	return nil, nil
}

func getValue(r *http.Request, key, defaultVal string) string {
	val := r.FormValue(key)
	log.Info("val", val)
	if val == "" {
		val = defaultVal
	}

	return val
}

func probeStart(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	ret := getValue(r, "result", "pass")

	if err := putProbe(0, ret, "start", 5); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func probeDrop(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe drop")
	ret := getValue(r, "result", "pass")
	timeout := r.FormValue("timeout")
	t, err := strconv.Atoi(timeout)
	if err != nil {
		util.RespHTTPErr(w, http.StatusBadRequest, err.Error())
	}
	time.Sleep(time.Duration(t) * time.Second)

	if err = putProbe(0, ret, "drop", 1); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	currKey := getCurrentKey()
	if err = getProbe(2, ret, currKey, "drop"+currKey); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe drop")
}

func probeRecoverPort(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe recover port")
	time.Sleep(5 * time.Second)
	ret := getValue(r, "result", "pass")

	if err := putProbe(0, ret, "recoverport", -1); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	currKey := getCurrentKey()
	if err := getProbe(0, ret, currKey, "recoverport"+currKey); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := getProbe(2, ret, currKey, "recoverport"+currKey); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe recover port")
}

func probePause(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe pause")
	ret := getValue(r, "result", "pass")

	if err := putProbe(0, ret, "pause", 1); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe pause")
}

func probeContinue(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe continue")
	time.Sleep(3 * time.Second)
	ret := getValue(r, "result", "pass")
	currKey := getCurrentKey()

	if err := getProbe(0, ret, currKey, "pause"+currKey); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := putProbe(0, ret, "continue", -1); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	currKey = getCurrentKey()
	if err := getProbe(2, ret, currKey, "continue"+currKey); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe continue")
}

func probeRestart(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe restart")
	time.Sleep(25 * time.Second)
	ret := getValue(r, "result", "pass")

	if err := putProbe(0, ret, "restart", -1); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	currKey := getCurrentKey()
	if err := getProbe(2, ret, currKey, "restart"+currKey); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe restart")
}

func getProbe(servID int, ret, key, val string) error {
	DB := makeDBClient(servAddrs[servID])
	kv, err := probePass(ret,
		func() (client.KeyValue, error) {
			return DB.Get(key)
		})
	if err != nil {
		return err
	}

	if string(kv.ValueBytes()) != val {
		err = errors.New("value unmatch")
		err = checkResult(err, ret)
	}

	log.Info("end: get, key:", key, "val:", val, "curr val:", string(kv.ValueBytes()))
	return err
}

func putProbe(servID int, ret, tag string, count int) error {
	var err error
	if count < 0 {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		count = int(rand.Int63n(100) + 1)
	}

	DB := makeDBClient(servAddrs[servID])
	for i := 0; i < count; i++ {
		key := generateKey()
		_, err = probePass(ret,
			func() (client.KeyValue, error) {
				return emptyKV, DB.Put(key, tag+key)
			})
		if err != nil {
			break
		}
	}

	log.Info("end: put, keys:", count)
	return err
}

func probeTest(w http.ResponseWriter, r *http.Request) {
	log.Info("probe")
	w.WriteHeader(http.StatusOK)
}
