package instance_agent

import (
	"fmt"
	"os/exec"
)

func DropPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptables -A OUTPUT -p tcp --dport %s -j drop", port)
	if _, err := exec.Command("/bin/sh", cmdStr).Output(); err != nil {
		return err
	}

	cmdStr = fmt.Sprintf("sudo iptables -A INPUT -p tcp -dport %s -j drop", port)
	_, err := exec.Command("/bin/sh", cmdStr).Output()

	return err
}

func RecoverPort(port string) error {
	cmdStr := fmt.Sprintf("sudo iptablse -D OUTPUT -p tcp --dport %s -j drop", port)
	if _, err := exec.Command("/bin/sh", cmdStr).Output(); err != nil {
		return err
	}

	cmdStr = fmt.Sprintf("sudo iptables -D INPUT -p tcp --dport %s -j drop", port)
	_, err := exec.Command("/bin/sh", cmdStr).Output()

	return err
}
