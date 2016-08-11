package logger

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/arsgo/color"
)

func (l *Logger) Info(content ...interface{}) {
	l.print(SLevel_Info, fmt.Sprint(content...))
}
func (l *Logger) Infof(format string, content ...interface{}) {
	l.Info(fmt.Sprintf(format, content...))
}

func (l *Logger) Debug(content ...interface{}) {
	l.print(SLevel_Debug, fmt.Sprint(content...))
}
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}
func (l *Logger) IFError(i bool, content ...interface{}) {
	if !i {
		return
	}
	l.Error(content...)
}
func (l *Logger) Error(content ...interface{}) {
	l.print(SLevel_Error, fmt.Sprint(content...))
}
func (l *Logger) IFErrorf(i bool, format string, a ...interface{}) {
	if !i {
		return
	}
	l.Errorf(format, a...)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}
func (l *Logger) Fatal(content ...interface{}) {
	l.print(SLevel_Fatal, fmt.Sprint(content...))
}
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}
func (l *Logger) Print(content ...interface{}) {
	l.Info(content...)
}
func (l *Logger) Printf(format string, a ...interface{}) {
	l.Infof(format, a...)
}

func (l *Logger) print(level string, content string) {
	defer l.recover()
	if strings.EqualFold(content, "") {
		return
	}
	events := l.getEvents(level, content)
	v := ""
	for _, event := range events {
		select {
		case l.DataChan <- event:
		default:
			v = ":已丢失"
		}
	}
	l.logPrint(level, content+v)
}

func (l *Logger) logPrint(level string, content string) {
	if !isDebug || !l.show {
		return
	}
	if levelMap[level] < levelMap[l.Level] {
		return
	}
	rcontext := content
	switch level {
	case SLevel_Debug:
		rcontext = color.MagentaString(content)
	case SLevel_Error:
		rcontext = color.YellowString(content)
	case SLevel_Fatal:
		rcontext = color.RedString(content)
	}
	log.Printf("[%s][%s]: %s", l.session, level[0:1], rcontext)

}

func (l *Logger) getEvents(level string, content string) (events []LoggerEvent) {
	events = make([]LoggerEvent, 0, 6)
	currentLevel := levelMap[level]
	if currentLevel < levelMap[l.Level] {
		return
	}

	event := LoggerEvent{Level: level, RLevel: level, Name: l.Name, Now: time.Now(), Content: content,
		Path: l.Config.Appender.Path, Session: l.session, Caller: getCaller(4)}
	events = append(events, event)

	//添加到info列表
	if currentLevel > ILevel_Info && currentLevel < ILevel_OFF {
		event := LoggerEvent{Level: SLevel_Info, RLevel: level, Name: l.Name, Now: time.Now(), Content: content,
			Path: l.Config.Appender.Path, Session: l.session, Caller: getCaller(4)}
		events = append(events, event)
	}

	/*for k, v := range levelMap {
		if v <= level && v <= logLevel && v > ILevel_OFF {
			event := &LoggerEvent{Level: k, Name: l.Name, Now: time.Now(), Content: content,
				Path: l.Config.Appender.Path}
			events[k] = event
		}
	}*/
	/*
		ILevel_ALL = iota
		ILevel_Debug
		ILevel_Info
		ILevel_Error
		ILevel_Fatal
		ILevel_OFF

	*/
	return
}
