package agent

import (
	"net/http"
	"os"
	"runtime"
	"sync"
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

func (s *TestNet) TestPortAlive(c *C) {
	addr := "127.0.0.1:" + port

	err := util.HTTPCall(util.JoinURL(addr, "test/net", ""), "POST", nil)
	c.Assert(err, IsNil)

	err = DropPort(port)
	c.Assert(err, IsNil)
	err = util.HTTPCall(util.JoinURL(addr, "test/net", ""), "POST", nil)
	c.Assert(err, NotNil)

	err = RecoverPort(port)
	c.Assert(err, IsNil)
	err = util.HTTPCall(util.JoinURL(addr, "test/net", ""), "POST", nil)
	c.Assert(err, IsNil)
}
