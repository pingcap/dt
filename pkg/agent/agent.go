package agent

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

const (
	registerIntervalTime = 1 //sec
)

type Agent struct {
	IP       string
	Addr     string
	CtrlAddr string

	l      net.Listener
	inst   *Instance
	exitCh chan string
}

func NewAgent(cfg *AgentConfig) (*Agent, error) {
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
		inst:     &Instance{state: instanceStateNone, logfile: f, cmd: exec.Command(cfg.DataDir)},
		exitCh:   make(chan string, 1)}, nil
}

func (a *Agent) Register() error {
	log.Debug("start: register")
	agentAttr := make(url.Values)
	agentAttr.Set("addr", a.Addr)

	return util.HttpCall(util.ApiUrl(a.CtrlAddr, "api/agent/register", agentAttr.Encode()), "POST", nil)
}

func (a *Agent) Start() error {
	for {
		if err := a.Register(); err != nil {
			log.Warning("register failed,errors.Trace(err):", errors.Trace(err))
			time.Sleep(registerIntervalTime * time.Millisecond)
		} else {
			break
		}
	}

	go runHttpServer(a)

	select {
	case msg := <-a.exitCh:
		if msg != "" {
			return errors.New(msg)
		}
	}

	return nil
}

func (a *Agent) Shutdown() error {
	log.Debug("start: shutdown")
	close(a.exitCh)

	return nil
}
