package logger

//ILogger 日志接口
type ILogger interface {
	Info(content ...interface{})
	Infof(format string, content ...interface{})
	Debug(content ...interface{})
	Debugf(format string, a ...interface{})
	Error(content ...interface{})
	Errorf(format string, a ...interface{})
	Fatal(content ...interface{})
	Fatalf(format string, a ...interface{})
	Printf(format string, a ...interface{})
	GetName()string
}
