package util

import (
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

func HttpCall(url, method string) error {
	log.Debug("start: httpCall, url:", url, "method:", method)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	data, _ := ioutil.ReadAll(resp.Body)
	if code == 200 {
		return nil
	}

	return errors.New(string(data))
}

func GetIpAndPort(addr string) (string, string, error) {
	strSlice := strings.Split(addr, ":")
	if len(strSlice) != 2 {
		return "", "", ErrAddrInvalid
	}

	return strSlice[0], strSlice[1], nil
}
