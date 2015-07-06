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

func (inst *Instance) execCmd(arg string) (*exec.Cmd, error) {
	cmd := exec.Command("sh", "-c", arg)
	cmd.Stdout = inst.logfile
	cmd.Stderr = inst.logfile

	return cmd, cmd.Run()
}

func (inst *Instance) Start(arg string) (err error) {
	log.Debug("start: startInstance, agent")
	if inst.cmd, err = inst.execCmd(arg); err != nil {
		return
	}

	inst.state = instanceStateStarted
	inst.pid = inst.cmd.Process.Pid

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
	if _, err := inst.execCmd(arg); err != nil {
		return errors.Trace(err)
	}
	inst.state = instanceStatePaused

	return nil
}

func (inst *Instance) ContinuePause() error {
	if inst.state != instanceStatePaused {
		return nil
	}

	arg := fmt.Sprintf(continueInstanceCmd+" %d", inst.pid)
	if _, err := inst.execCmd(arg); err != nil {
		return errors.Trace(err)
	}
	inst.state = instanceStateStarted

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

	return nil
}

func (inst *Instance) BackupData(args ...string) error {
	arg := fmt.Sprintf(backupInstanceDataCmd+" %s", inst.dataDir)
	_, err := inst.execCmd(arg)

	return errors.Trace(err)
}

func (inst *Instance) CleanUpData(args ...string) error {
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
