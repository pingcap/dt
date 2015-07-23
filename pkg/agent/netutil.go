package agent

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/pingcap/dt/pkg/util"
)

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
