package logger

import "time"

const (
	ILevel_ALL = iota
	ILevel_Debug
	ILevel_Info
	ILevel_Error
	ILevel_Fatal
	ILevel_OFF
)
const (
	SLevel_OFF   = "Off"
	SLevel_Info  = "Info"
	SLevel_Error = "Error"
	SLevel_Fatal = "Fatal"
	SLevel_Debug = "Debug"
	SLevel_ALL   = "All"
)

type LoggerAppender struct {
	Type  string
	Level string
	Path  string
}
type LoggerLayout struct {
	Level   int
	Content string
}
type LoggerConfig struct {
	Name     string
	Appender *LoggerAppender
}

//LoggerEvent 日志
type LoggerEvent struct {
	Level   string
	Now     time.Time
	Name    string
	Content string
	Path    string
	Session string
	Caller  string
}
