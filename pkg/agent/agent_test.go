package agent

import (
	"runtime"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type TestAgent struct{}

var _ = Suite(&TestAgent{})

// TODO:
func (s *TestAgent) SetUpSuite(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// TODO:
func (s *TestAgent) TearDownSuite(c *C) {}
