package controller

import (
	"github.com/BurntSushi/toml"
	"github.com/ngaut/log"
)

//controller config
type CtrlCfg struct {
	Title         string
	Addr          string
	DataDir       string                  `toml:"data_dir"`
	InstanceCount int                     `toml:"instance_count"`
	InstanceInfos map[string]InstanceInfo `toml:"instance_info"`
	Cmds          []TestCmd               `toml:"test_cmd"`
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

func GetCfg(file string) (cfg *CtrlCfg, err error) {
	log.Debug("start: getCtrlCfg")
	_, err = toml.DecodeFile(file, &cfg)
	log.Info(cfg)

	return
}
