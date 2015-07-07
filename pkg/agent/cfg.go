package agent

import (
	"github.com/BurntSushi/toml"
)

//agent config
type Config struct {
	IP       string
	Port     string
	CtrlAddr string `toml:"ctrl_addr"`
	DataDir  string `toml:"data_dir"`
}

func LoadConfig(file string) (cfg *Config, err error) {
	_, err = toml.DecodeFile(file, &cfg)

	return
}
