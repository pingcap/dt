package agent

import (
	"fmt"
	"net/url"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

type Agent struct {
	IP       string
	Addr     string
	CtrlAddr string
	DataDir  string

	inst   *Instance
	exitCh chan error
}

func NewAgent(cfg *Config) (*Agent, error) {
	f, err := util.CreateLog(cfg.DataDir, "agent_cmd")
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Agent{
		IP:       cfg.IP,
		Addr:     fmt.Sprintf("%s:%s", cfg.IP, cfg.Port),
		CtrlAddr: cfg.CtrlAddr,
		DataDir:  cfg.DataDir,
		inst:     NewInstance(f),
		exitCh:   make(chan error, 1)}, nil
}

func (a *Agent) heartbeat() error {
	agentAttr := make(url.Values)
	agentAttr.Set("addr", a.Addr)

	return util.HTTPCall(util.JoinURL(a.CtrlAddr, "api/agent/register", agentAttr.Encode()), "POST", nil)
}

func (a *Agent) Register() error {
	log.Debug("start: register")
	for {
		if err := a.heartbeat(); err != nil {
			log.Warning("register failed, err:", err)
			time.Sleep(util.HeartbeatIntervalSec)
		} else {
			break
		}
	}

	return nil
}

func (a *Agent) Heartbeat() {
	log.Debug("start: heartbeat")
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := a.heartbeat(); err != nil {
				log.Warning("heartbeat failed, err - ", err)
				break
			}
			log.Debug("heartbeat")
		case <-a.exitCh:
			return
		}
	}
}

func (a *Agent) Start() error {
	go runHTTPServer(a)
	a.Register()
	go a.Heartbeat()

	select {
	case err := <-a.exitCh:
		time.Sleep(1 * time.Second)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (a *Agent) BackupData(path string) error {
	arg := fmt.Sprintf("%s %s %s", backupInstanceDataCmd, a.DataDir, path)
	if _, err := util.ExecCmd(arg, a.inst.logfile); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (a *Agent) CleanupData() error {
	arg := fmt.Sprintf("%s %s", cleanUpInstanceDataCmd, a.DataDir)
	if _, err := util.ExecCmd(arg, a.inst.logfile); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (a *Agent) Shutdown() error {
	log.Debug("start: shutdown")
	close(a.exitCh)

	return nil
}
