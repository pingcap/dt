package util

import (
	"github.com/BurntSushi/toml"
)

const (
	CmdStartInstance = "start"
)

//controller config
type TestCfg struct {
	Title        string
	Attr         *BaseAttr
	InstanceInfo map[string]InstanceInfo
	Cmds         []*TestCmd `toml: "testCmd"`
}

type BaseAttr struct {
	Addr          string
	DataDir       string
	InstanceCount int
}

type InstanceInfo struct {
	Count int
}

type TestCmd struct {
	Name      string `toml: "cmd"`
	Dir       string
	Args      string
	Probe     string
	Instances []string
}

//instance_agent config
type AgentCfg struct {
	Attr *AgentAttr
}

type AgentAttr struct {
	Ip       string
	Port     string
	CtrlAddr string
	DataDir  string
}

func GetCtrlCfg(file string) (cfg *TestCfg, err error) {
	_, err = toml.DecodeFile(file, &cfg)

	return
}

func GetAgentCfg(file string) (cfg *AgentCfg, err error) {
	_, err = toml.DecodeFile(file, &cfg)

	return
}
