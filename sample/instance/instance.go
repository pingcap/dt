package main

import (
	"flag"
	"net"
	"os"
	"runtime"

	"github.com/ngaut/log"
)

var (
	logDir = flag.String("log-dir", "./dt/instance", "the directory to store log")
	addr   = flag.String("addr", ":54300", "http listen addr")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	os.MkdirAll(*logDir, 0755)
	log.SetLevelByString("debug")
	log.SetOutputByName(*logDir + "/inst.log")

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err = l.Accept()
		if err != nil {
			log.Warning("accept err - ", err)
		}
		// TODO: do sth
	}
}
