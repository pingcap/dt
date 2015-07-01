package cmd

import (
	"flag"
	"log"
	"runtime"

	ctrl "testingframe/pkg/controller"
	agent "testingframe/pkg/instance_agent"
)

var (
	role = flag.String("role", "instance_agent", "start the specified process")
	path = flag.String("path", "/tmp", "the path of data directory")
	addr = flag.String("addr", "127.0.0.1:54321", "http listen address")
	cfg  = flag.String("cfg", "cmd/cfg.toml", "configure file name")
)

type Server interface {
	Start(arg string) error
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	var s Server
	var err error

	switch *role {
	case "instanc_agent":
		s, err = agent.NewInstanceAgent(*path, *addr)
		if err != nil {
			log.Fatalln(err)
		}
	case "controller":
		s = ctrl.NewController(*path, *addr)
	}

	log.Fatalln(s.Start(*cfg))
}
