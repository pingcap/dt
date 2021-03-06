package agent

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/util"
)

const (
	instanceStateUninit   = "uninit"
	instanceStateStarted  = "started"
	instanceStateStopped  = "stopped"
	instanceStatePaused   = "paused"
	instanceStateContinue = "continue"

	stopInstanceCmd        = "kill -KILL"
	pauseInstanceCmd       = "kill -STOP"
	continueInstanceCmd    = "kill -CONT"
	backupInstanceDataCmd  = "cp -r"
	cleanUpInstanceDataCmd = "rm -r"
)

type Instance struct {
	pid     int
	state   string
	logfile *os.File
	cmd     *exec.Cmd
}

func NewInstance(f *os.File) *Instance {
	return &Instance{state: instanceStateUninit, logfile: f}
}

// TODO: used for checking results
func ps() string {
	cmd := exec.Command("sh", "-c", "ps -aux|grep cockroach")
	output, _ := cmd.Output()

	return string(output)
}

// TODO: used for checking results
func listIPTables() string {
	cmd := exec.Command("sh", "-c", "sudo iptables -L")
	output, _ := cmd.Output()

	return string(output)
}

func (inst *Instance) Start(arg, name string) error {
	var err error
	pidFile := fmt.Sprintf("%s.out", util.GetGUID(name))
	isNohup := strings.Contains(arg, "nohup")

	if isNohup {
		arg = fmt.Sprintf("%s echo $! > %s", arg, pidFile)
	}
	if inst.cmd, err = util.ExecCmd(arg, inst.logfile); err != nil {
		return errors.Trace(err)
	}

	defer log.Warning("start out:", ps(), "cmd:", arg, "pid:", inst.pid)

	inst.state = instanceStateStarted
	if isNohup {
		buf, err := util.ReadFile(pidFile)
		if err == nil {
			inst.pid, err = strconv.Atoi(string(buf))
		}
		if err != nil {
			return errors.Trace(err)
		}

		return err
	}
	inst.pid = inst.cmd.Process.Pid

	return err
}

func (inst *Instance) Set(arg string) error {
	log.Debug("start: set, arg:", arg)
	var err error
	if _, err = util.ExecCmd(arg, inst.logfile); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (inst *Instance) Restart(arg, name string) error {
	if err := inst.Stop(); err != nil {
		return errors.Trace(err)
	}

	return inst.Start(arg, name)
}

func (inst *Instance) Pause() error {
	if inst.state != instanceStateStarted {
		return nil
	}

	args := fmt.Sprintf("%s %d", pauseInstanceCmd, inst.pid)
	cmd, err := util.ExecCmd(args, inst.logfile)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()
	inst.state = instanceStatePaused

	log.Warning("pause out:", ps())

	return nil
}

func (inst *Instance) Continue() error {
	if inst.state != instanceStatePaused {
		return nil
	}

	args := fmt.Sprintf("%s %d", continueInstanceCmd, inst.pid)
	if _, err := util.ExecCmd(args, inst.logfile); err != nil {
		return errors.Trace(err)
	}
	inst.state = instanceStateStarted

	log.Warning("continue out:", ps())

	return nil
}

func (inst *Instance) Stop() error {
	if inst.state != instanceStateStarted {
		return nil
	}

	args := fmt.Sprintf("%s %d", stopInstanceCmd, inst.pid)
	cmd, err := util.ExecCmd(args, inst.logfile)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Process.Wait()

	inst.state = instanceStateStopped
	log.Warning("stop out:", ps())

	return nil
}

func (inst *Instance) DropPort(port string) error {
	return DropPort(port)
}

func (inst *Instance) DropPkg(chain, port string, percent int) error {
	return DropPkg(chain, port, percent)
}

func (inst *Instance) LimitSpeed(chain, port, unit string, pkgs int) error {
	return LimitSpeed(chain, port, unit, pkgs)
}

func (inst *Instance) RecoverPort(port string) error {
	return RecoverPort(port)
}
