package agent

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/pingcap/dt/pkg/util"
)

func DropPkg(chain, port string, percent int) error {
	cmdStr := fmt.Sprintf("sudo iptables -A %s -p tcp --dport %s -m statistic --mode random --probability %f -j DROP",
		chain, port, float32(percent)/100)
	cmd, err := util.ExecCmd(cmdStr, os.Stdout)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()

	return nil
}

func LimitSpeed(chain, port, unit string, pkgs int) error {
	cmdStr := fmt.Sprintf("sudo iptables -A %s -p tcp --dport %s -m limit --limit %d/%s -j ACCEPT",
		chain, port, pkgs, unit)
	cmd, err := util.ExecCmd(cmdStr, os.Stdout)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()

	cmdStr = fmt.Sprintf("sudo iptables -A %s -p tcp --dport %s -j DROP", chain, port)
	cmd, err = util.ExecCmd(cmdStr, os.Stdout)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()

	return nil
}

// TODO:
func NetDelay(port string, delay int) error { return nil }

func DropPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptables -A OUTPUT -p tcp --dport %s -j DROP -w", port)
	cmd, err := util.ExecCmd(cmdStr, os.Stdout)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()

	cmdStr = fmt.Sprintf("sudo iptables -A INPUT -p tcp --dport %s -j DROP -w", port)
	if cmd, err = util.ExecCmd(cmdStr, os.Stdout); err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()

	return nil
}

func RecoverPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptables -D OUTPUT -p tcp --dport %s -j DROP", port)
	cmd, err := util.ExecCmd(cmdStr, os.Stdout)
	if err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()

	cmdStr = fmt.Sprintf("sudo iptables -D INPUT -p tcp --dport %s -j DROP", port)
	if cmd, err = util.ExecCmd(cmdStr, os.Stdout); err != nil {
		return errors.Trace(err)
	}
	cmd.Wait()

	return nil
}
