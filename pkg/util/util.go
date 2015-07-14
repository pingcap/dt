package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
)

func ApiUrl(addr, action, query string) string {
	if query != "" {
		return fmt.Sprintf("http://%s/%s?%s", addr, action, query)
	} else {
		return fmt.Sprintf("http://%s/%s", addr, action)
	}
}

func HTTPCall(url, method string, data interface{}) error {
	log.Debug("start: httpCall, url:", url, "method:", method)
	rw := &bytes.Buffer{}
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			return errors.Trace(err)
		}
		rw.Write(buf)
	}
	req, err := http.NewRequest(method, url, rw)
	if err != nil {
		return errors.Trace(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode/100 != 2 {
		return errors.Errorf("error code: %s, msg: %s", resp.StatusCode, string(msg))
	}

	return nil
}

func RespHTTPErr(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	if msg == "" {
		msg = http.StatusText(code)
	}
	io.WriteString(w, msg)
}

func ExecCmd(arg string, w io.Writer) (*exec.Cmd, error) {
	cmd := exec.Command("sh", "-c", arg)
	cmd.Stdout = w
	cmd.Stderr = w

	return cmd, cmd.Start()
}

func GetGUID(key string) string {
	t := time.Now().UnixNano()

	return fmt.Sprintf("%d-%s", t, key)
}

func ReadFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Trace(err)
	}

	buf := bytes.Trim(b, "\n")

	return buf, nil
}

func CheckIsEmpty(strs ...string) bool {
	return Contains("", strs)
}

func Contains(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}
