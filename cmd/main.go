package main

import (
	"flag"
	"runtime"

	"github.com/ngaut/log"
	ctrl "github.com/pingcap/dt/pkg/controller"
	agent "github.com/pingcap/dt/pkg/instance_agent"
	"github.com/pingcap/dt/pkg/util"
)

var (
	role    = flag.String("role", "instance_agent", "start the specified process")
	cfgPath = flag.String("cfg", "cmd/cfg.toml", "configure file name")
	level   = flag.String("loglevel", "debug", "set log level: info, warn, error, debug [default: info]")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetLevelByString(*level)

	switch *role {
	case "instanc_agent":
		cfg, err := util.GetAgentCfg(*cfgPath)
		if err != nil {
			log.Fatal(err)
		}
		s, err := agent.NewInstanceAgent(cfg.Attr.DataDir, cfg.Attr.Ip, cfg.Attr.Port,
			cfg.Attr.CtrlAddr)
		if err == nil {
			err = s.Start()
		}
		log.Fatal(err)
	case "controller":
		cfg, err := util.GetCtrlCfg(*cfgPath)
		if err != nil {
			log.Fatal(err)
		}
		s := ctrl.NewController(cfg.Attr.DataDir, cfg.Attr.Addr)
		log.Fatal(s.Start(cfg))
	}
}
