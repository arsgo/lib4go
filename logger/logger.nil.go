package logger

import "fmt"

type NilLogger struct {
	show bool
}

func (l *NilLogger) Info(content ...interface{}) {
	l.print(SLevel_Info, fmt.Sprint(content...))
}
func (l *NilLogger) Infof(format string, content ...interface{}) {
	l.Info(fmt.Sprintf(format, content...))
}

func (l *NilLogger) Debug(content ...interface{}) {
	l.print(SLevel_Debug, fmt.Sprint(content...))
}
func (l *NilLogger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}
func (l *NilLogger) IFError(i bool, content ...interface{}) {
	if !i {
		return
	}
	l.Error(content...)
}
func (l *NilLogger) Error(content ...interface{}) {
	l.print(SLevel_Error, fmt.Sprint(content...))
}
func (l *NilLogger) IFErrorf(i bool, format string, a ...interface{}) {
	if !i {
		return
	}
	l.Errorf(format, a...)
}

func (l *NilLogger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}
func (l *NilLogger) Fatal(content ...interface{}) {
	l.print(SLevel_Fatal, fmt.Sprint(content...))
}
func (l *NilLogger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}
func (l *NilLogger) Print(content ...interface{}) {
	l.Info(content...)
}
func (l *NilLogger) Printf(format string, a ...interface{}) {
	l.Infof(format, a...)
}
func (l *NilLogger) print(level string, content string) {
	fmt.Println(content)
}
func (l *NilLogger) GetName() string {
	return "nil"
}
func (l *NilLogger) Show(b bool) {
	l.show = b
}
