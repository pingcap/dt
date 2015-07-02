package cmd

import (
	"flag"
	"log"
	"runtime"

	ctrl "github.com/pingcap/dt/pkg/controller"
	agent "github.com/pingcap/dt/pkg/instance_agent"
	"github.com/pingcap/dt/pkg/util"
)

var (
	role    = flag.String("role", "instance_agent", "start the specified process")
	cfgPath = flag.String("cfg", "cmd/cfg.toml", "configure file name")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error

	switch *role {
	case "instanc_agent":
		cfg, err := util.GetAgentCfg(*cfgPath)
		if err != nil {
			log.Fatalln(err)
		}
		s, err := agent.NewInstanceAgent(cfg.Attr.DataDir, cfg.Attr.Ip, cfg.Attr.Port,
			cfg.Attr.CtrlAddr)
		if err == nil {
			err = s.Start()
		}

	case "controller":
		cfg, err := util.GetCtrlCfg(*cfgPath)
		if err != nil {
			log.Fatalln(err)
		}
		s := ctrl.NewController(cfg.Attr.DataDir, cfg.Attr.Addr)
		err = s.Start(cfg)
	}

	log.Fatalln(err)
}
