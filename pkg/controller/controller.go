package controller

import (
	"fmt"
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
	// argument split
	argSplit = "/"
)

var (
	errUnmatchedCfgInfo      = errors.New("unmatched config info")
	errAgentRegisterTimeout  = errors.New("register timeout")
	errAgentHeartbeatTimeout = errors.New("heartbeat timeout")
	errUnmatchedTestCmd      = errors.New("unmatched test cmd kind")
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
		return nil, errors.Trace(errUnmatchedCfgInfo)
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
	log.Debug("start: getAgentAddrs, count:", ctrl.getAgentsCount())
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
	interval := 5 * util.HeartbeatIntervalSec
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
		log.Warning("handle cmd failed, err - ", err)
		err = ctrl.HandleFailure()
		return errors.Trace(err)
	}

	return nil
}

func (ctrl *Controller) HandleFailure() error {
	instances := make([]string, 0, len(ctrl.agents))
	for name, _ := range ctrl.agents {
		instances = append(instances, name)
	}

	// TODO: restart all instance
	//	restartCmd := &TestCmd{Name: util.TestCmdRestart, Args, Probe, Instances: instances}
	//	if err := ctrl.HandleCmd(restartCmd); err != nil {
	//		return errors.Trace(err)
	//	}

	return nil
}

func (ctrl *Controller) HandleCmd(cmd *TestCmd) error {
	if (cmd.Name != util.TestCmdSleep && len(cmd.Instances) <= 0) ||
		(cmd.Name == util.TestCmdSleep && len(cmd.Instances) > 0) {
		return errors.Trace(errUnmatchedCfgInfo)
	}

	if cmd.Name == util.TestCmdSleep {
		if err := DoCmd(cmd, nil, ""); err != nil {
			return errors.Trace(err)
		}
	}

	for _, inst := range cmd.Instances {
		agent, ok := ctrl.agents[inst]
		if !ok {
			return errors.Trace(errUnmatchedCfgInfo)
		}
		if err := DoCmd(cmd, agent, inst); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func DoCmd(cmd *TestCmd, agent *client.Agent, inst string) error {
	var err error
	log.Debug("start: docmd, cmd:", cmd.Name)
	defer log.Debug("end: docmd")

	switch strings.ToLower(cmd.Name) {
	case util.TestCmdStart:
		err = agent.StartInstance(cmd.Args, inst, cmd.Probe)
	case util.TestCmdInit:
		err = agent.SetInstance(cmd.Args, cmd.Probe)
	case util.TestCmdRestart:
		err = agent.RestartInstance(cmd.Args, inst, cmd.Probe)
	case util.TestCmdPause:
		err = agent.PauseInstance(cmd.Probe)
	case util.TestCmdContinue:
		err = agent.ContinueInstance(cmd.Probe)
	case util.TestCmdStop:
		err = agent.StopInstance(cmd.Probe)
	case util.TestCmdDropPort:
		err = agent.DropPortInstance(cmd.Args, cmd.Probe)
	case util.TestCmdDropPkg:
		args := strings.Split(cmd.Args, argSplit)
		if len(args) != 3 {
			err = errUnmatchedCfgInfo
			break
		}
		err = agent.DropPkgInstance(args[0], args[1], args[2], cmd.Probe)
	case util.TestCmdLimitSeep:
		args := strings.Split(cmd.Args, argSplit)
		if len(args) != 4 {
			err = errUnmatchedCfgInfo
			break
		}
		err = agent.LimitSpeedInstance(args[0], args[1], args[2], args[3], cmd.Probe)
	case util.TestCmdRecoverPort:
		err = agent.RecoverPortInstance(cmd.Args, cmd.Probe)
	case util.TestCmdShutdownAgent:
		err = agent.Shutdown()
	case util.TestCmdBackupData:
		err = agent.BackupInstanceData(cmd.Dir)
	case util.TestCmdCleanUpData:
		err = agent.CleanUpInstanceData()
	case util.TestCmdSleep:
		t, err := strconv.Atoi(cmd.Args)
		if err != nil {
			return errors.Trace(err)
		}
		time.Sleep(time.Duration(t) * time.Second)
	default:
		err = errUnmatchedTestCmd
	}

	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
