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
		inst:     NewInstance(f),
		exitCh:   make(chan error, 1)}, nil
}

func (a *Agent) heartbeat() error {
	agentAttr := make(url.Values)
	agentAttr.Set("addr", a.Addr)

	return util.HTTPCall(util.ApiUrl(a.CtrlAddr, "api/agent/register", agentAttr.Encode()), "POST", nil)
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

func (a *Agent) Heartbeat() error {
	log.Debug("start: heartbeat")
	for {
		if err := a.heartbeat(); err != nil {
			log.Warning("hb failed, errors.Trace(err):", errors.Trace(err))
		}
		time.Sleep(3 * time.Second)
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

func (a *Agent) Shutdown() error {
	log.Debug("start: shutdown")
	close(a.exitCh)

	return nil
}
