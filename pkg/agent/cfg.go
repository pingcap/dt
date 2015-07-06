package agent

import (
	"github.com/BurntSushi/toml"
)

//agent config
type AgentConfig struct {
	IP       string
	Port     string
	CtrlAddr string `toml:"ctrl_addr"`
	DataDir  string `toml:"data_dir"`
}

func LoadConfig(file string) (cfg *AgentConfig, err error) {
	_, err = toml.DecodeFile(file, &cfg)

	return
}
