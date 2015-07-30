package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"
	"github.com/ngaut/log"
)

func JoinURL(addr, action, query string) string {
	if query != "" {
		return fmt.Sprintf("http://%s/%s?%s", addr, action, query)
	} else {
		return fmt.Sprintf("http://%s/%s", addr, action)
	}
}

func HTTPCall(url, method string, data interface{}) error {
	log.Debug("start: httpCall, url:", url, "method:", method)
	defer log.Debug("end: httpCall")
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

	req.Header.Set("Connection", "close")
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
