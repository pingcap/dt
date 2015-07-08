package main

import (
	"flag"
	"runtime"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/agent"
	ctrl "github.com/pingcap/dt/pkg/controller"
)

var (
	role    = flag.String("role", "agent", "start the specified process: controller, agent [default: agent]")
	cfgPath = flag.String("cfg", "etc/agent_cfg.toml", "configure file name")
	level   = flag.String("loglevel", "debug", "set log level: info, warn, error, debug [default: debug]")
)

type Server interface {
	Start() error
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetLevelByString(*level)

	var s Server
	var err error

	switch *role {
	case "agent":
		cfg, err := agent.LoadConfig(*cfgPath)
		log.Debug(cfg)
		if err != nil {
			log.Fatal(errors.ErrorStack(err))
		}
		s, err = agent.NewAgent(cfg)
	case "controller":
		cfg, err := ctrl.LoadConfig(*cfgPath)
		log.Debug(cfg)
		if err != nil {
			log.Fatal(errors.ErrorStack(err))
		}
		s, err = ctrl.NewController(cfg)
	}

	if err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
	if err = s.Start(); err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
}
