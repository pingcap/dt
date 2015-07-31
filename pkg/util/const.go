package util

import (
	"time"
)

// controller test cmd name
const (
	TestCmdStart         = "start"
	TestCmdInit          = "init"
	TestCmdRestart       = "restart"
	TestCmdPause         = "pause"
	TestCmdContinue      = "continue"
	TestCmdStop          = "stop"
	TestCmdSleep         = "sleep"
	TestCmdDropPkg       = "droppkg"
	TestCmdLimitSeep     = "limitspeed"
	TestCmdDropPort      = "dropport"
	TestCmdRecoverPort   = "recoverport"
	TestCmdCleanUpData   = "cleanup"
	TestCmdBackupData    = "backup"
	TestCmdShutdownAgent = "shutdown"
)

// heartbeat interval
const (
	HeartbeatIntervalSec = time.Second
)

//HTTP
const (
	HTTPClientTimeout = 10 * time.Second
)
