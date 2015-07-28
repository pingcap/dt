package agent

import (
	"net/http"
	"os"
	"time"

	"github.com/pingcap/dt/pkg/util"
	. "gopkg.in/check.v1"
)

const (
	port = "9876"
)

func listenAndServe(c *C) {
	http.HandleFunc("/test/port", testPortAlive)

	err := http.ListenAndServe(":"+port, nil)
	c.Assert(err, IsNil)
}

func testPortAlive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *TestAgent) Test(c *C) {
	addr := "127.0.0.1:" + port

	_, err := util.ExecCmd("sudo iptables -F", os.Stdout)
	c.Assert(err, IsNil)
	go listenAndServe(c)
	time.Sleep(time.Second)

	err = util.HTTPCall(util.JoinURL(addr, "test/port", ""), "POST", nil)
	c.Assert(err, IsNil)

	err = DropPort(port)
	c.Assert(err, IsNil)
	err = util.HTTPCall(util.JoinURL(addr, "test/port", ""), "POST", nil)
	c.Assert(err, NotNil)

	err = RecoverPort(port)
	c.Assert(err, IsNil)
	err = util.HTTPCall(util.JoinURL(addr, "test/port", ""), "POST", nil)
	c.Assert(err, IsNil)
}
