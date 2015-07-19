package agent

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/pingcap/dt/pkg/util"
)

func DropPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptables -A OUTPUT -p tcp --dport %s -j DROP -w", port)
	_, err := util.ExecCmd(cmdStr, os.Stdout)
	if err != nil {
		return errors.Trace(err)
	}

	cmdStr = fmt.Sprintf("sudo iptables -A INPUT -p tcp --dport %s -j DROP -w", port)
	if _, err = util.ExecCmd(cmdStr, os.Stdout); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func RecoverPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptables -D OUTPUT -p tcp --dport %s -j DROP", port)
	if _, err := util.ExecCmd(cmdStr, os.Stdout); err != nil {
		return errors.Trace(err)
	}

	cmdStr = fmt.Sprintf("sudo iptables -D INPUT -p tcp --dport %s -j DROP", port)
	if _, err := util.ExecCmd(cmdStr, os.Stdout); err != nil {
		return errors.Trace(err)
	}

	return nil
}
