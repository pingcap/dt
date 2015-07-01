package controller

import (
	"errors"
	"time"

	client "testingframe/pkg/instance_agent/client"
	"testingframe/pkg/util"
)

const (
	agentInfoChanSize    = 20
	agentRegisterTimeout = 60
)

var (
	ErrCfgInfoUnmatch       = errors.New("unmath config info")
	ErrAgentRegisterTimeout = errors.New("register timeout")
)

type Controller struct {
	Addr       string
	DataDir    string
	AgentInfos []*AgentInfo

	agentCount  int
	agents      map[string]*client.Agent
	cmds        []*util.TestCmd
	agentInfoCh chan *AgentInfo
}

type AgentInfo struct {
	dir  string
	ip   string
	addr string
}

func NewController(dataDir, addr string) *Controller {
	return &Controller{
		Addr:        addr,
		DataDir:     dataDir,
		agentInfoCh: make(chan *AgentInfo, agentInfoChanSize)}
}

func getAgentInfos(count int, infoCh chan *AgentInfo) ([]*AgentInfo, error) {
	agentInfos := make([]*AgentInfo, count)
	timeout := time.After(agentRegisterTimeout * time.Second)

	for {
		select {
		case info := <-infoCh:
			agentInfos = append(agentInfos, info)
		case <-timeout:
			return nil, ErrAgentRegisterTimeout
		}
		if agentInfos[count-1] != nil {
			break
		}
	}

	return agentInfos, nil
}

func (ctrl *Controller) Init(cfgFile string) (err error) {
	cfg, err := util.GetCfg(cfgFile)
	if err != nil {
		return
	}

	instanceCount := 0
	for _, inst := range cfg.InstanceInfo {
		instanceCount += inst.Count
	}
	if cfg.Attr.InstanceCount != len(cfg.Cmds) || cfg.Attr.InstanceCount != instanceCount {
		return ErrCfgInfoUnmatch
	}

	ctrl.cmds = cfg.Cmds
	ctrl.agentCount = cfg.Attr.InstanceCount
	ctrl.agents = make(map[string]*client.Agent, ctrl.agentCount)
	ctrl.AgentInfos, err = getAgentInfos(ctrl.agentCount, ctrl.agentInfoCh)

	return
}

func (ctrl *Controller) Start(cfgFile string) error {
	if err := ctrl.Init(cfgFile); err != nil {
		return err
	}

	for _, cmd := range ctrl.cmds {
		if err := ctrl.ExecCmd(cmd); err != nil {
			// TODO: deal with failure
		}
	}

	return nil
}

// TODO: implement
func (ctrl *Controller) ExecCmd(cmd *util.TestCmd) error {
	panic("ExecCmd hasn't implemented")

	return nil
}
