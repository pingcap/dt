package controller

import (
	"errors"
	"strconv"
	"time"

	"github.com/ngaut/log"
	client "github.com/pingcap/dt/pkg/instance_agent/client"
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
	cmds        []util.TestCmd
	agentInfoCh chan string
}

func NewController(dataDir, addr string) *Controller {
	return &Controller{
		Addr:        addr,
		DataDir:     dataDir,
		agentInfoCh: make(chan string, agentInfoChanSize)}
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

func (ctrl *Controller) Init(cfg *util.CtrlCfg) (err error) {
	log.Debug("start: init")
	instanceCount := 0
	for _, inst := range cfg.InstanceInfos {
		instanceCount += inst.Count
	}
	if cfg.Attr.InstanceCount != instanceCount {
		return ErrCfgInfoUnmatch
	}

	ctrl.Addr = cfg.Attr.Addr
	ctrl.cmds = cfg.Cmds
	ctrl.agentCount = cfg.Attr.InstanceCount
	ctrl.agents = make(map[string]*client.Agent, ctrl.agentCount)
	instanceCount = 1
	for kind, inst := range cfg.InstanceInfos {
		for i := 0; i < inst.Count; i++ {
			ctrl.agents[kind+strconv.Itoa(instanceCount)] = &client.Agent{}
			instanceCount++
		}
	}

	return
}

func (ctrl *Controller) Start(cfgFile *util.CtrlCfg) error {
	if err := ctrl.Init(cfgFile); err != nil {
		return err
	}

	go runHttpServer(ctrl.Addr, ctrl)
	if err := ctrl.getAgentAddrs(); err != nil {
		return err
	}

	log.Info(ctrl.cmds)
	for _, cmd := range ctrl.cmds {
		if err := ctrl.HandleCmd(cmd); err != nil {
			// TODO: deal with failure
		}
	}

	return nil
}

// name, dir, args, probe, instances
func (ctrl *Controller) HandleCmd(cmd util.TestCmd) error {
	log.Debug("start: handlecmd")
	switch cmd.Name {
	case util.TestCmdStart:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].StartInstance(cmd.Args, cmd.Probe); err != nil {
				return err
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
	default:
		return ErrTestCmdUnmatch
	}

	return nil
}
