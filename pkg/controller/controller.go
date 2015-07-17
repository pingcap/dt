package controller

import (
	"fmt"
	"net"
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
	errCfgInfoUnmatch        = errors.New("unmatch config info")
	errAgentRegisterTimeout  = errors.New("register timeout")
	errAgentHeartbeatTimeout = errors.New("heartbeat timeout")
	errTestCmdUnmatch        = errors.New("test cmd kind unmatch")
)

type Controller struct {
	Addr    string
	DataDir string

	agents      map[string]*client.Agent
	cmds        []*TestCmd
	agentInfoCh chan string
	exitCh      chan error
}

func NewController(cfg *Config) (*Controller, error) {
	ctrl := &Controller{
		Addr:        cfg.Addr,
		DataDir:     cfg.DataDir,
		cmds:        cfg.Cmds,
		exitCh:      make(chan error, 1),
		agentInfoCh: make(chan string, agentInfoChanSize*3)}

	instanceCount := 0
	for _, inst := range cfg.InstanceInfos {
		instanceCount += inst.Count
	}
	if cfg.InstanceCount != instanceCount {
		return nil, errors.Trace(errCfgInfoUnmatch)
	}

	ctrl.agents = make(map[string]*client.Agent, cfg.InstanceCount)
	index := 1
	for kind, inst := range cfg.InstanceInfos {
		for i := 0; i < inst.Count; i++ {
			agent := fmt.Sprintf("%s%d", kind, index)
			ctrl.agents[agent] = &client.Agent{}
			index++
		}
	}

	return ctrl, nil
}

func (ctrl *Controller) getAgentsCount() int {
	return len(ctrl.agents)
}

func (ctrl *Controller) getAgentAddrs() error {
	log.Debug("start: getAgentAddrs")
	agentAddrs := make([]string, ctrl.getAgentsCount())

	i := 0
	lastAddr := ctrl.getAgentsCount() - 1
	timeout := time.After(agentRegisterTimeout * time.Second)
	for {
		select {
		case addr := <-ctrl.agentInfoCh:
			if util.Contains(addr, agentAddrs) {
				break
			}
			agentAddrs[i] = addr
			i++
		case <-timeout:
			return errors.Trace(errAgentRegisterTimeout)
		}
		if agentAddrs[lastAddr] != "" {
			break
		}
	}

	i = 0
	var err error
	for _, agent := range ctrl.agents {
		agent.Addr = agentAddrs[i]
		agent.LastHeartbeat = time.Now()
		if agent.Ip, _, err = net.SplitHostPort(agentAddrs[i]); err != nil {
			return errors.Trace(err)
		}
		i++
	}

	return nil
}

func (ctrl *Controller) checkAlive() {
	interval := 3 * util.HeartbeatIntervalSec
	t := time.NewTicker(interval)
	defer t.Stop()

	setHeartbeat := func(addr string) {
		for _, agent := range ctrl.agents {
			if agent.Addr == addr {
				agent.LastHeartbeat = time.Now()
				break
			}
		}
	}

	checkTimeout := func() {
		for _, agent := range ctrl.agents {
			if time.Now().Sub(agent.LastHeartbeat) > interval {
				ctrl.exitCh <- errAgentHeartbeatTimeout
				return
			}
		}
	}

	for {
		select {
		case addr := <-ctrl.agentInfoCh:
			setHeartbeat(addr)
		case <-t.C:
			checkTimeout()
		}
	}
}

func (ctrl *Controller) Start() error {
	go runHTTPServer(ctrl.Addr, ctrl)
	if err := ctrl.getAgentAddrs(); err != nil {
		return errors.Trace(err)
	}
	go ctrl.checkAlive()

	var err error
	for _, cmd := range ctrl.cmds {
		if len(ctrl.exitCh) > 0 {
			// TODO: clean up data
			return errors.Trace(<-ctrl.exitCh)
		}
		if err = ctrl.HandleCmd(cmd); err == nil {
			continue
		}
		if err = ctrl.HandleFailure(); err != nil {
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
	for _, inst := range cmd.Instances {
		if err := ctrl.DoCmd(cmd, inst); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (ctrl *Controller) DoCmd(cmd *TestCmd, inst string) error {
	log.Debug("start: docmd, cmd:", cmd.Name)
	switch strings.ToLower(cmd.Name) {
	case util.TestCmdStart:
		err := ctrl.agents[inst].StartInstance(cmd.Args, inst, cmd.Dir, cmd.Probe)
		if err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdInit:
		err := ctrl.agents[inst].SetInstance(cmd.Args, cmd.Probe)
		if err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdRestart:
		err := ctrl.agents[inst].RestartInstance(cmd.Args, inst, cmd.Dir, cmd.Probe)
		if err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdPause:
		if err := ctrl.agents[inst].PauseInstance(cmd.Probe); err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdContinue:
		if err := ctrl.agents[inst].ContinueInstance(cmd.Probe); err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdStop:
		if err := ctrl.agents[inst].StopInstance(); err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdDropPort:
		if err := ctrl.agents[inst].DropPortInstance(cmd.Args, cmd.Probe); err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdRecoverPort:
		if err := ctrl.agents[inst].RecoverPortInstance(cmd.Args, cmd.Probe); err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdShutdownAgent:
		if err := ctrl.agents[inst].Shutdown(); err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdBackupData:
		log.Info("backup, dir:", cmd.Dir)
		if err := ctrl.agents[inst].BackupInstanceData(cmd.Dir); err != nil {
			return errors.Trace(err)
		}
	case util.TestCmdCleanUpData:
		if err := ctrl.agents[inst].CleanUpInstanceData(); err != nil {
			return errors.Trace(err)
		}
	default:
		return errors.Trace(errTestCmdUnmatch)
	}

	return nil
}
