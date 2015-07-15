package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/dt/pkg/agent"
	ctrl "github.com/pingcap/dt/pkg/controller"
	"github.com/pingcap/dt/pkg/util"
)

var (
	role     = flag.String("role", "agent", "start the specified process: controller, agent [default: agent]")
	cfgPath  = flag.String("cfg", "etc/agent_cfg.toml", "configure file name")
	logLevel = flag.String("log-level", "debug", "set log level: info, warn, error, debug [default: debug]")
	logDir   = flag.String("log-dir", "./", "the directory to store log")
)

type Server interface {
	Start() error
}

func setLogInfo(level, logDir, file string) error {
	if _, err := util.CreateLog(logDir, file); err != nil {
		return errors.Trace(err)
	}

	path := fmt.Sprintf("%s/%s.log", logDir, file)
	log.SetLevelByString(level)
	log.SetOutputByName(path)

	return nil
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := setLogInfo(*logLevel, *logDir, *role); err != nil {
		log.Fatal(errors.ErrorStack)
	}

	var s Server
	var err error

	switch *role {
	case "agent":
		cfg, err := agent.LoadConfig(*cfgPath)
		if err != nil {
			log.Fatal(errors.ErrorStack(err))
		}
		s, err = agent.NewAgent(cfg)
	case "controller":
		cfg, err := ctrl.LoadConfig(*cfgPath)
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
