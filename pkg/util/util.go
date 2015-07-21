package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/juju/errors"
)

func ExecCmd(arg string, w io.Writer) (*exec.Cmd, error) {
	cmd := exec.Command("sh", "-c", arg)
	cmd.Stdout = w
	cmd.Stderr = w

	return cmd, cmd.Start()
}

func GetGUID(key string) string {
	t := time.Now().UnixNano()

	return fmt.Sprintf("%d-%s", t, key)
}

func ReadFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Trace(err)
	}

	buf := bytes.Trim(b, "\n")

	return buf, nil
}

func CreateLog(dir, file string) (*os.File, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, errors.Trace(err)
	}

	path := fmt.Sprintf("%s/%s.log", dir, file)
	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return f, nil
}

func CheckIsEmpty(strs ...string) bool {
	return Contains("", strs)
}

func Contains(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}
