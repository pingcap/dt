package client

import (
	"net/url"
	"time"

	"github.com/pingcap/dt/pkg/util"
)

type Agent struct {
	Ip            string
	Addr          string
	LastHeartbeat time.Time
}

func NewAgent(dir, addr, ip string) (*Agent, error) {
	return &Agent{Ip: ip, Addr: addr}, nil
}

func (a *Agent) SetInstance(cmd, probe string) error {
	attr := make(url.Values)
	attr.Set("cmd", cmd)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/set", attr.Encode()), "POST", nil)
}

func (a *Agent) StartInstance(cmd, instName, probe string) error {
	attr := make(url.Values)
	attr.Set("cmd", cmd)
	attr.Set("probe", probe)
	attr.Set("name", instName)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/start", attr.Encode()), "POST", nil)
}

func (a *Agent) RestartInstance(cmd, instName, probe string) error {
	attr := make(url.Values)
	attr.Set("cmd", cmd)
	attr.Set("probe", probe)
	attr.Set("name", instName)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/restart", attr.Encode()), "POST", nil)
}

func (a *Agent) PauseInstance(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/pause", attr.Encode()), "POST", nil)
}

func (a *Agent) ContinueInstance(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/continue", attr.Encode()), "POST", nil)
}

func (a *Agent) BackupInstanceData(dir string) error {
	attr := make(url.Values)
	attr.Set("dir", dir)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/backupdata", attr.Encode()), "POST", nil)
}

func (a *Agent) CleanUpInstanceData() error {
	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/cleanupdata", ""), "POST", nil)
}

func (a *Agent) StopInstance(probe string) error {
	attr := make(url.Values)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/stop", attr.Encode()), "POST", nil)
}

func (a *Agent) DropPkgInstance(port, chain, percent, probe string) error {
	attr := make(url.Values)
	attr.Set("port", port)
	attr.Set("chain", chain)
	attr.Set("percent", percent)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/droppkg", attr.Encode()), "POST", nil)
}

func (a *Agent) LimitSpeedInstance(port, chain, unit, pkgs, probe string) error {
	attr := make(url.Values)
	attr.Set("port", port)
	attr.Set("chain", chain)
	attr.Set("unit", unit)
	attr.Set("pkgs", pkgs)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/limitspeed", attr.Encode()), "POST", nil)
}

func (a *Agent) DropPortInstance(port, probe string) error {
	attr := make(url.Values)
	attr.Set("port", port)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/dropport", attr.Encode()), "POST", nil)
}

func (a *Agent) RecoverPortInstance(port, probe string) error {
	attr := make(url.Values)
	attr.Set("port", port)
	attr.Set("probe", probe)

	return util.HTTPCall(util.JoinURL(a.Addr, "api/instance/recoverport", attr.Encode()), "POST", nil)
}

func (a *Agent) Shutdown() error {
	return util.HTTPCall(util.JoinURL(a.Addr, "api/agent/shutdown", ""), "POST", nil)
}
