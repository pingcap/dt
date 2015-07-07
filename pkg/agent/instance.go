package agent

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/juju/errors"
	"github.com/ngaut/log"
)

const (
	instanceStateNone     = "uninitialized"
	instanceStateStarted  = "started"
	instanceStateStopped  = "stopped"
	instanceStatePaused   = "paused"
	instanceStateContinue = "continue"

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

func (inst *Instance) execCmd(arg string) (*exec.Cmd, error) {
	cmd := exec.Command("sh", "-c", arg)
	cmd.Stdout = inst.logfile
	cmd.Stderr = inst.logfile

	return cmd, cmd.Start()
}

// TODO: used for testing
func ps() string {
	cmd := exec.Command("sh", "-c", "ps -ejf|grep test")
	output, _ := cmd.Output()

	return string(output)
}

func (inst *Instance) Start(arg string) (err error) {
	log.Debug("start: startInstance, agent")
	if inst.cmd, err = inst.execCmd(arg); err != nil {
		return
	}

	inst.state = instanceStateStarted
	inst.pid = inst.cmd.Process.Pid

	log.Warning("start out:", ps())

	return
}

func (inst *Instance) Restart(arg string) error {
	if err := inst.Stop(); err != nil {
		return errors.Trace(err)
	}

	return inst.Start(arg)
}

func (inst *Instance) Pause() error {
	if inst.state != instanceStateStarted {
		return nil
	}

	arg := fmt.Sprintf(pauseInstanceCmd+" %d", inst.pid)
	_, err := inst.execCmd(arg)
	if err != nil {
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

	arg := fmt.Sprintf(continueInstanceCmd+" %d", inst.pid)
	if _, err := inst.execCmd(arg); err != nil {
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

	if err := inst.cmd.Process.Kill(); err != nil {
		return errors.Trace(err)
	}
	if _, err := inst.cmd.Process.Wait(); err != nil {
		return errors.Trace(err)
	}

	inst.state = instanceStateStopped
	log.Warning("stop out:", ps())

	return nil
}

func (inst *Instance) BackupData(path string) error {
	arg := fmt.Sprintf(backupInstanceDataCmd+" %s %s", inst.dataDir, path)
	_, err := inst.execCmd(arg)

	return errors.Trace(err)
}

func (inst *Instance) CleanUpData() error {
	// TODO: clean up intance internal state
	arg := fmt.Sprintf(cleanUpInstanceDataCmd+" %s", inst.dataDir)
	_, err := inst.execCmd(arg)

	return errors.Trace(err)
}

//  TODO: ports may be more than one
func (inst *Instance) DropPort(port string) error {
	return DropPort(port)
}

func (inst *Instance) RecoverPort(port string) error {
	return RecoverPort(port)
}
