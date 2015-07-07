package agent

import (
	"fmt"
	"net/url"
	"os"
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
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return nil, errors.Trace(err)
	}

	f, err := os.Create(cfg.DataDir + "/agent.log")
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

func (a *Agent) Register() error {
	log.Debug("start: register")
	agentAttr := make(url.Values)
	agentAttr.Set("addr", a.Addr)

	return util.HTTPCall(util.ApiUrl(a.CtrlAddr, "api/agent/register", agentAttr.Encode()), "POST", nil)
}

func (a *Agent) Start() error {
	go runHTTPServer(a)
	for {
		if err := a.Register(); err != nil {
			log.Warning("register failed,errors.Trace(err):", errors.Trace(err))
			time.Sleep(1 * time.Millisecond)
		} else {
			break
		}
	}

	select {
	case err := <-a.exitCh:
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
