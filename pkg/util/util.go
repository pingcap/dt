package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/juju/errors"
	"github.com/ngaut/log"
)

var ErrAddrInvalid = errors.New("invalid addr")

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

func GetIpAndPort(addr string) (string, string, error) {
	strSlice := strings.Split(addr, ":")
	if len(strSlice) != 2 {
		return "", "", errors.Trace(ErrAddrInvalid)
	}

	return strSlice[0], strSlice[1], nil
}

func WriteHTTPResponse(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func WriteHTTPError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, msg)
}
