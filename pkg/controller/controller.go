package controller

import (
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/agent/client"
	"github.com/pingcap/dt/pkg/util"
)

const (
	agentInfoChanSize    = 20
	agentRegisterTimeout = 60
)

var (
	ErrCfgInfoUnmatch       = errors.New("unmath config info")
	ErrAgentRegisterTimeout = errors.New("register timeout")
	ErrTestCmdUnmatch       = errors.New("test cmd kind unmath")
)

type Controller struct {
	Addr    string
	DataDir string

	agentCount  int
	agents      map[string]*client.Agent
	cmds        []TestCmd
	agentInfoCh chan string
}

func NewController(cfg *CtrlConfig) (*Controller, error) {
	ctrl := &Controller{
		Addr:        cfg.Addr,
		DataDir:     cfg.DataDir,
		agentInfoCh: make(chan string, agentInfoChanSize)}

	instanceCount := 0
	for _, inst := range cfg.InstanceInfos {
		instanceCount += inst.Count
	}
	if cfg.InstanceCount != instanceCount {
		return nil, errors.Trace(ErrCfgInfoUnmatch)
	}

	ctrl.Addr = cfg.Addr
	ctrl.cmds = cfg.Cmds
	ctrl.agentCount = cfg.InstanceCount
	ctrl.agents = make(map[string]*client.Agent, ctrl.agentCount)
	instanceCount = 1
	for kind, inst := range cfg.InstanceInfos {
		for i := 0; i < inst.Count; i++ {
			ctrl.agents[kind+strconv.Itoa(instanceCount)] = &client.Agent{}
			instanceCount++
		}
	}

	return ctrl, nil
}

func (ctrl *Controller) getAgentAddrs() (err error) {
	log.Debug("start: getAgentAddrs")
	agentAddrs := make([]string, ctrl.agentCount)
	timeout := time.After(agentRegisterTimeout * time.Second)
	i := 0

	for {
		select {
		case info := <-ctrl.agentInfoCh:
			agentAddrs[i] = info
			i++
		case <-timeout:
			return ErrAgentRegisterTimeout
		}
		if agentAddrs[ctrl.agentCount-1] != "" {
			break
		}
	}

	i = 0
	for _, agent := range ctrl.agents {
		agent.Addr = agentAddrs[i]
		if agent.Ip, _, err = util.GetIpAndPort(agentAddrs[i]); err != nil {

			return
		}
		i++
	}

	return
}

func (ctrl *Controller) Start() error {
	go runHttpServer(ctrl.Addr, ctrl)
	if err := ctrl.getAgentAddrs(); err != nil {
		return errors.Trace(err)
	}

	for _, cmd := range ctrl.cmds {
		log.Info(cmd.Name)
		if err := ctrl.HandleCmd(cmd); err != nil {
			// TODO: deal with failure
			return errors.Trace(err)
		}
	}

	return nil
}

func (ctrl *Controller) HandleCmd(cmd TestCmd) error {
	log.Debug("start: handlecmd")
	switch strings.ToLower(cmd.Name) {
	case util.TestCmdStart:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].StartInstance(cmd.Args, cmd.Probe); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdRestart:
		// TODO: implement
		panic("ExecCmd hasn't implemented")
	case util.TestCmdPause:
		// TODO: implement
		panic("ExecCmd hasn't implemented")
	case util.TestCmdContinue:
		// TODO: implement
		panic("ExecCmd hasn't implemented")
	case util.TestCmdStop:
		// TODO: implement
		panic("ExecCmd hasn't implemented")
	case util.TestCmdDropPort:
		// TODO: implement
		panic("ExecCmd hasn't implemented")
	case util.TestCmdRecoverPort:
		// TODO: implement
		panic("ExecCmd hasn't implemented")
	case util.TestCmdShutdownAgent:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].Shutdown(); err != nil {
				return errors.Trace(err)
			}
		}
	default:
		return errors.Trace(ErrTestCmdUnmatch)
	}

	return nil
}
