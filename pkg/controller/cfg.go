package util

import (
	"github.com/BurntSushi/toml"
	"github.com/ngaut/log"
)

const (
	CmdStartInstance = "start"
)

//controller config
type CtrlCfg struct {
	Title         string
	Attr          BaseAttr                `toml:"baseAttr"`
	InstanceInfos map[string]InstanceInfo `toml:"instanceInfo"`
	Cmds          []TestCmd               `toml:"testCmd"`
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
	Name      string `toml:"cmd"`
	Dir       string
	Args      string
	Probe     string
	Instances []string
}

//instance_agent config
type AgentCfg struct {
	Attr AgentAttr
}

type AgentAttr struct {
	Ip       string
	Port     string
	CtrlAddr string
	DataDir  string
}

func GetCtrlCfg(file string) (cfg *CtrlCfg, err error) {
	log.Debug("start: getCtrlCfg")
	_, err = toml.DecodeFile(file, &cfg)
	log.Info(cfg)

	return
}

func GetAgentCfg(file string) (cfg *AgentCfg, err error) {
	log.Debug("start: getAgentCfg")
	_, err = toml.DecodeFile(file, &cfg)
	log.Info(cfg)

	return
}
