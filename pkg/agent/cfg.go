package agent

import (
	"github.com/BurntSushi/toml"
	"github.com/ngaut/log"
)

//agent config
type AgentCfg struct {
	Ip       string
	Port     string
	CtrlAddr string `toml:"ctrl_addr"`
	DataDir  string `toml:"data_dir"`
}

func GetCfg(file string) (cfg *AgentCfg, err error) {
	log.Debug("start: getAgentCfg")
	_, err = toml.DecodeFile(file, &cfg)
	log.Info(cfg)

	return
}
