package util

// instance_agent api url
const (
	ActionStartInstance       = "api/start_instance"
	ActionRestartInstance     = "api/restart_instance"
	ActionPauseInstance       = "api/pause_instance"
	ActionContinueInstance    = "api/continue_instance"
	ActionStopInstance        = "api/stop_instance"
	ActionDropPortInstance    = "api/dropport_instance"
	ActionRecoverPortInstance = "api/Recoverport_instance"
	ActionBackupInstanceData  = "api/backupdata_instance"
	ActionCleanUpInstanceData = "api/cleanUpdata_instance"
	ActionShutdown            = "api/Shutdown"
)

// ctroller api url
const (
	ActionRegisterAgent = "api/register_agent"
)

// ctroller test cmd name
const (
	TestCmdStart       = "start"
	TestCmdRestart     = "restart"
	TestCmdPause       = "pause"
	TestCmdContinue    = "continue"
	TestCmdStop        = "stop"
	TestCmdDropPort    = "dropport"
	TestCmdRecoverPort = "recoverport"
)
