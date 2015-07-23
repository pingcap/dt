package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
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
	timeoutFlag = true
	passResult  = "pass"
)

var (
	keyGlobal = 1
	emptyKV   = client.KeyValue{}
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

func checkResult(err error, result string, flag bool) error {
	if flag == timeoutFlag {
		if result == "pass" {
			return errors.New("timeout")
		}
		return nil
	}

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
		return &kv, checkResult(ret, isPass, !timeoutFlag)
	case <-timeout:
		return &kv, checkResult(nil, isPass, timeoutFlag)
	}

	return nil, nil
}

func getResult(ret string) string {
	if ret == "" {
		ret = passResult
	}

	return ret
}

func probeStart(w http.ResponseWriter, r *http.Request) {
	key := generateKey()
	time.Sleep(5 * time.Second)

	DB := makeDBClient(servAddrs[0])
	if err := DB.Put(key, "start"+key); err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func probeDropPort(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe drop port")
	ret := getResult(r.FormValue("result"))
	timeout := r.FormValue("timeout")
	t, err := strconv.Atoi(timeout)
	if err != nil {
		util.RespHTTPErr(w, http.StatusBadRequest, err.Error())
	}
	key := generateKey()
	val := "dorpport" + key

	time.Sleep(time.Duration(t) * time.Second)

	DB := makeDBClient(servAddrs[0])
	_, err = probePass(ret,
		func() (client.KeyValue, error) {
			return emptyKV, DB.Put(key, val)
		})
	if err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	DB = makeDBClient(servAddrs[2])
	kv, err := probePass(ret,
		func() (client.KeyValue, error) {
			return DB.Get(key)
		})
	if err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if string(kv.ValueBytes()) != val {
		err = errors.New("value unmatch")
		if err = checkResult(err, ret, !timeoutFlag); err != nil {
			util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe drop port")
}

func probeRecoverPort(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe recover port")
	ret := getResult(r.FormValue("result"))
	key := generateKey()
	time.Sleep(5 * time.Second)

	DB := makeDBClient(servAddrs[0])
	_, err := probePass(ret,
		func() (client.KeyValue, error) {
			return emptyKV, DB.Put(key, "recoverport"+key)
		})
	if err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe recover port")
}

func probePause(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe pause")
	ret := getResult(r.FormValue("result"))
	key := generateKey()

	DB := makeDBClient(servAddrs[0])
	_, err := probePass(ret,
		func() (client.KeyValue, error) {
			return emptyKV, DB.Put(key, "pause"+key)
		})
	if err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe pause")
}

func probeContinue(w http.ResponseWriter, r *http.Request) {
	log.Debug("start: probe continue")
	ret := getResult(r.FormValue("result"))
	key := generateKey()

	DB := makeDBClient(servAddrs[0])
	_, err := probePass(ret,
		func() (client.KeyValue, error) {
			return emptyKV, DB.Put(key, "continue"+key)
		})
	if err != nil {
		util.RespHTTPErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Debug("end: probe continue")
}

func probeTest(w http.ResponseWriter, r *http.Request) {
	log.Info("probe")
	w.WriteHeader(http.StatusOK)
}
