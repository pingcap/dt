package controller

import (
	"net"
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
	agentRegisterTimeout = 300
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
	cmds        []*TestCmd
	agentInfoCh chan string
}

func NewController(cfg *Config) (*Controller, error) {
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
		case addr := <-ctrl.agentInfoCh:
			if util.CheckIsExist(addr, agentAddrs) {
				break
			}
			agentAddrs[i] = addr
			i++
		case <-timeout:
			return errors.Trace(ErrAgentRegisterTimeout)
		}
		if agentAddrs[ctrl.agentCount-1] != "" {
			break
		}
	}

	i = 0
	for _, agent := range ctrl.agents {
		agent.Addr = agentAddrs[i]
		if agent.Ip, _, err = net.SplitHostPort(agentAddrs[i]); err != nil {
			return errors.Trace(err)
		}
		i++
	}

	return
}

func (ctrl *Controller) Start() error {
	go runHTTPServer(ctrl.Addr, ctrl)
	if err := ctrl.getAgentAddrs(); err != nil {
		return errors.Trace(err)
	}

	for _, cmd := range ctrl.cmds {
		if err := ctrl.HandleCmd(cmd); err == nil {
			continue
		}
		if err := ctrl.HandleFailure(); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (ctrl *Controller) HandleFailure() error {
	instances := make([]string, 0, len(ctrl.agents))
	for name, _ := range ctrl.agents {
		instances = append(instances, name)
	}

	cleanCmd := &TestCmd{Name: util.TestCmdCleanUpData, Instances: instances}
	if err := ctrl.HandleCmd(cleanCmd); err != nil {
		return errors.Trace(err)
	}

	// TODO: restart all instance
	//	restartCmd := &TestCmd{Name: util.TestCmdRestart, Args, Probe, Instances: instances}
	//	if err := ctrl.HandleCmd(restartCmd); err != nil {
	//		return errors.Trace(err)
	//	}

	return nil
}

func (ctrl *Controller) HandleCmd(cmd *TestCmd) error {
	log.Debug("start: handlecmd, cmd:", cmd.Name)
	switch strings.ToLower(cmd.Name) {
	case util.TestCmdStart:
		for _, inst := range cmd.Instances {
			err := ctrl.agents[inst].StartInstance(cmd.Args, inst, cmd.Dir, cmd.Probe)
			if err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdRestart:
		for _, inst := range cmd.Instances {
			err := ctrl.agents[inst].RestartInstance(cmd.Args, inst, cmd.Dir, cmd.Probe)
			if err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdPause:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].PauseInstance(cmd.Probe); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdContinue:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].ContinueInstance(cmd.Probe); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdStop:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].StopInstance(cmd.Probe); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdDropPort:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].DropPortInstance(cmd.Args, cmd.Probe); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdRecoverPort:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].RecoverPortInstance(cmd.Args, cmd.Probe); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdShutdownAgent:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].Shutdown(); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdBackupData:
		for _, inst := range cmd.Instances {
			log.Info("backup, dir:", cmd.Dir)
			if err := ctrl.agents[inst].BackupInstanceData(cmd.Dir); err != nil {
				return errors.Trace(err)
			}
		}
	case util.TestCmdCleanUpData:
		for _, inst := range cmd.Instances {
			if err := ctrl.agents[inst].CleanUpInstanceData(); err != nil {
				return errors.Trace(err)
			}
		}
	default:
		return errors.Trace(ErrTestCmdUnmatch)
	}

	return nil
}
