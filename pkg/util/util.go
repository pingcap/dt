package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

func HttpCall(url, method string, data interface{}) error {
	log.Debug("start: httpCall, url:", url, "method:", method)
	rw := &bytes.Buffer{}
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			return err
		}
		rw.Write(buf)
	}
	req, err := http.NewRequest(method, url, rw)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	msg, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 == 2 {
		return nil
	}

	return errors.New(fmt.Sprintf("error code: %s, msg: %s", resp.StatusCode, string(msg)))
}

func GetIpAndPort(addr string) (string, string, error) {
	strSlice := strings.Split(addr, ":")
	if len(strSlice) != 2 {
		return "", "", ErrAddrInvalid
	}

	return strSlice[0], strSlice[1], nil
}
