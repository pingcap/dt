package agent

import (
	"net/http"
	"os"
	"runtime"
	//	"sync"
	"time"

	"github.com/pingcap/dt/pkg/util"
	. "gopkg.in/check.v1"
)

const (
	port = "9876"
)

type TestNet struct{}

var _ = Suite(&TestNet{})

func (s *TestNet) SetUpSuite(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go listenAndServe(c)
	time.Sleep(time.Second)
}

func (s *TestNet) SetUpTest(c *C) {
	_, err := util.ExecCmd("sudo iptables -F", os.Stdout)
	c.Assert(err, IsNil)
}

func (s *TestNet) TearDownSuite(c *C) {
	_, err := util.ExecCmd("sudo iptables -F", os.Stdout)
	c.Assert(err, IsNil)
}

func listenAndServe(c *C) {
	http.HandleFunc("/test/net", testNet)

	err := http.ListenAndServe(":"+port, nil)
	c.Assert(err, IsNil)
}

func testNet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// TODO: test or code has problems, still need to think it over
//func (s *TestNet) TestLimitSpeed(c *C) {
//	total := 20
//	pkgs := 5
//	failed := 0
//	unit := "second"
//	addr := "127.0.0.1:" + port
//
//	err := LimitSpeed("INPUT", port, unit, pkgs)
//	c.Assert(err, IsNil)
//
//	wg := sync.WaitGroup{}
//	data := make([]byte, 1500)
//	for i := 0; i < total; i++ {
//		wg.Add(1)
//		go func(f *int) {
//			err := util.HTTPCall(util.JoinURL(addr, "test/net", ""), "POST", data)
//			if err != nil {
//				*f++
//			}
//			wg.Done()
//		}(&failed)
//	}
//	wg.Wait()
//	c.Assert(failed, Equals, total-pkgs)
//}
//

func (s *TestNet) TestDropPkg(c *C) {
	failed := 0
	count := 4
	percent := 25
	addr := "127.0.0.1:" + port

	err := DropPkg("INPUT", port, percent)
	c.Assert(err, IsNil)
	for i := 0; i < count; i++ {
		err = util.HTTPCall(util.JoinURL(addr, "test/net", ""), "POST", nil)
		if err != nil {
			failed++
		}
	}
	c.Assert(failed, Equals, count*percent/100)
}

func (s *TestNet) TestPortAlive(c *C) {
	addr := "127.0.0.1:" + port

	err := DropPort(port)
	c.Assert(err, IsNil)
	err = util.HTTPCall(util.JoinURL(addr, "test/net", ""), "POST", nil)
	c.Assert(err, NotNil)

	err = RecoverPort(port)
	c.Assert(err, IsNil)
	err = util.HTTPCall(util.JoinURL(addr, "test/net", ""), "POST", nil)
	c.Assert(err, IsNil)
}
