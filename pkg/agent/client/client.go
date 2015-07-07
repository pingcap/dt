package client

import (
	"net/url"

	"github.com/juju/errors"
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

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/start", attr.Encode()), "POST", nil)
}

func (a *Agent) RestarInstance(cmd, probe string) error {
	attr := make(url.Values)
	attr.Set("cmd", cmd)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/restart", attr.Encode()), "POST", nil)
}

func (a *Agent) PauseInstance(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/pause", attr.Encode()), "POST", nil)
}

func (a *Agent) ContinueInstance(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/continue", attr.Encode()), "POST", nil)
}

func (a *Agent) BackupInstanceData(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/backupdata", attr.Encode()), "POST", nil)
}

func (a *Agent) CleanUpInstanceData(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/cleanupdata", attr.Encode()), "POST", nil)
}

func (a *Agent) StopInstance(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/stop", attr.Encode()), "POST", nil)
}

func (a *Agent) DropPortInstance(port, probe string) error {
	attr := make(url.Values)
	attr.Set("port", port)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/dropport", attr.Encode()), "POST", nil)
}

func (a *Agent) RecoverPortInstance(port, probe string) error {
	attr := make(url.Values)
	attr.Set("port", port)
	attr.Set("probe", probe)

	return util.HTTPCall(util.ApiUrl(a.Addr, "api/instance/recoverport", attr.Encode()), "POST", nil)
}

func (a *Agent) Shutdown() error {
	return util.HTTPCall(util.ApiUrl(a.Addr, "api/agent/shutdown", ""), "POST", nil)
}
