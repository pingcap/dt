package agent

import (
	"fmt"
	"os/exec"

	"github.com/juju/errors"
)

func DropPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptables -A OUTPUT -p tcp --dport %s -j drop", port)
	if _, err := exec.Command("/bin/sh", cmdStr).Output(); err != nil {
		return errors.Trace(err)
	}

	cmdStr = fmt.Sprintf("sudo iptables -A INPUT -p tcp -dport %s -j drop", port)
	_, err := exec.Command("/bin/sh", cmdStr).Output()

	return errors.Trace(err)
}

func RecoverPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptablse -D OUTPUT -p tcp --dport %s -j drop", port)
	if _, err := exec.Command("/bin/sh", cmdStr).Output(); err != nil {
		return errors.Trace(err)
	}

	cmdStr = fmt.Sprintf("sudo iptables -D INPUT -p tcp --dport %s -j drop", port)
	_, err := exec.Command("/bin/sh", cmdStr).Output()

	return errors.Trace(err)
}
