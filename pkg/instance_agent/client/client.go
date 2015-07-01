package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"testingframe/pkg/util"
)

var (
	ErrResposeCodeUnmath = errors.New("respose code unmath")
)

type Agent interface {
	StartInstance(args ...string) error
	RestarInstance(args ...string) error
	PauseInstance() error
	ConitnueInstace() error
	BackupInstanceData(args ...string) error
	CleanUpInstanceData(args ...string) error
	StopInstance() error
	DropPortInstance(port string) error
	RecoverPortInstance(port string) error
	Shutdown() error
}

type agent struct {
	dir  string
	ip   string
	addr string
}

func post(data []byte, url string) (*http.Response, error) {
	c := &http.Client{}
	buff := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func NewAgent(dir, addr, ip string) (Agent, error) {
	return &agent{dir: dir, ip: ip, addr: addr}, nil
}

func (a *agent) StartInstance(args ...string) error {
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	resp, err := post(b, util.UrlStartInstance)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		//TODO: add log
		return ErrResposeCodeUnmath
	}

	return nil
}

// TODO: implement
func (a *agent) RestarInstance(args ...string) error      { return nil }
func (a *agent) PauseInstance() error                     { return nil }
func (a *agent) ConitnueInstace() error                   { return nil }
func (a *agent) BackupInstanceData(args ...string) error  { return nil }
func (a *agent) CleanUpInstanceData(args ...string) error { return nil }
func (a *agent) StopInstance() error                      { return nil }
func (a *agent) DropPortInstance(port string) error       { return nil }
func (a *agent) RecoverPortInstance(port string) error    { return nil }
func (a *agent) Shutdown() error                          { return nil }
