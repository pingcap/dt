package controller

import (
	"github.com/BurntSushi/toml"
)

//controller config
type Config struct {
	Addr          string
	DataDir       string                  `toml:"data_dir"`
	InstanceCount int                     `toml:"instance_count"`
	InstanceInfos map[string]InstanceInfo `toml:"instance"`
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

func LoadConfig(file string) (cfg *Config, err error) {
	_, err = toml.DecodeFile(file, &cfg)

	return
}
