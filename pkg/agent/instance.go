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
	instanceStateNone     = "uninitialized"
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
	dataDir string
	logfile *os.File
	cmd     *exec.Cmd
}

func NewInstance(f *os.File) *Instance {
	return &Instance{state: instanceStateNone, logfile: f}
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

func (inst *Instance) Start(arg, name string) (err error) {
	log.Debug("start: startInstance, agent")
	pidFile := fmt.Sprintf("%s.out", util.GetGUID(name))
	isNohup := strings.Contains(arg, "nohup")

	if isNohup {
		arg = fmt.Sprintf("%s echo $! > %s", arg, pidFile)
	}
	if inst.cmd, err = util.ExecCmd(arg, inst.logfile); err != nil {
		return errors.Trace(err)
	}

	log.Warning("start out:", ps(), "cmd:", arg)

	inst.state = instanceStateStarted
	if isNohup {
		buf, err := util.ReadFile(pidFile)
		if err == nil {
			inst.pid, err = strconv.Atoi(string(buf))
		}
		if err != nil {
			return errors.Trace(err)
		}

		log.Warning("nohup, pid:", inst.pid)
		return err
	}
	inst.pid = inst.cmd.Process.Pid
	log.Warning("pid:", inst.pid)

	return
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

	arg := fmt.Sprintf("%s %d", pauseInstanceCmd, inst.pid)
	if _, err := util.ExecCmd(arg, inst.logfile); err != nil {
		return errors.Trace(err)
	}
	inst.state = instanceStatePaused

	log.Warning("pause out:", ps())

	return nil
}

func (inst *Instance) Continue() error {
	if inst.state != instanceStatePaused {
		return nil
	}

	arg := fmt.Sprintf("%s %d", continueInstanceCmd, inst.pid)
	if _, err := util.ExecCmd(arg, inst.logfile); err != nil {
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

	arg := fmt.Sprintf("%s %d", stopInstanceCmd, inst.pid)
	cmd, err := util.ExecCmd(arg, inst.logfile)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Process.Wait()

	inst.state = instanceStateStopped
	log.Warning("stop out:", ps())

	return nil
}

func (inst *Instance) BackupData(path string) error {
	arg := fmt.Sprintf("%s %s %s", backupInstanceDataCmd, inst.dataDir, path)
	_, err := util.ExecCmd(arg, inst.logfile)

	return errors.Trace(err)
}

func (inst *Instance) CleanUpData() error {
	// TODO: clean up intance internal state
	arg := fmt.Sprintf("%s %s", cleanUpInstanceDataCmd, inst.dataDir)
	_, err := util.ExecCmd(arg, inst.logfile)

	return errors.Trace(err)
}

func (inst *Instance) DropPort(port string) error {
	if err := DropPort(port); err != nil {
		return errors.Trace(err)
	}

	//log.Warning("drop port out:", listIPTables())

	return nil
}

func (inst *Instance) RecoverPort(port string) error {
	if err := RecoverPort(port); err != nil {
		return errors.Trace(err)
	}

	//log.Warning("recover port out:", listIPTables())

	return nil
}
