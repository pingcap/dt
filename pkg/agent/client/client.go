package client

import (
	"errors"
	"net/url"

	"github.com/pingcap/dt/pkg/util"
)

var (
	ErrResposeCodeUnmath = errors.New("respose code unmath")
)

type Agent struct {
	dir  string
	Ip   string
	Addr string
}

func NewAgent(dir, addr, ip string) (*Agent, error) {
	return &Agent{dir: dir, Ip: ip, Addr: addr}, nil
}

func (a *Agent) StartInstance(cmd, probe string) error {
	attr := make(url.Values)
	attr.Set("cmd", cmd)
	attr.Set("probe", probe)
	attr.Set("dir", a.dir)

	return util.HttpCall(util.ApiUrl(a.Addr, "api/instance/start", attr.Encode()), "POST", nil)
}

// TODO: implement
func (a *Agent) RestarInstance(args ...string) error      { return nil }
func (a *Agent) PauseInstance() error                     { return nil }
func (a *Agent) ConitnueInstace() error                   { return nil }
func (a *Agent) BackupInstanceData(args ...string) error  { return nil }
func (a *Agent) CleanUpInstanceData(args ...string) error { return nil }
func (a *Agent) StopInstance() error                      { return nil }
func (a *Agent) DropPortInstance(port string) error       { return nil }
func (a *Agent) RecoverPortInstance(port string) error    { return nil }

func (a *Agent) Shutdown() error {
	return util.HttpCall(util.ApiUrl(a.Addr, "api/agent/shutdown", ""), "POST", nil)
}
