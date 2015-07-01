package util

import (
	"github.com/BurntSushi/toml"
)

const (
	CmdStartInstance = "start"
)

type TestCfg struct {
	Titel        string
	Attr         baseAttr
	InstanceInfo map[string]instanceInfo
	Cmds         []*TestCmd `toml: "testCmd"`
}

type baseAttr struct {
	InstanceCount int
}

type instanceInfo struct {
	Count int
}

type TestCmd struct {
	Name      string `toml: "cmd"`
	Dir       string
	Args      string
	Probe     string
	Instances []string
}

func GetCfg(file string) (cfg *TestCfg, err error) {
	_, err = toml.DecodeFile(file, &cfg)

	return
}
