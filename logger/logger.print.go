package logger

import (
	"fmt"
	"log"
	"strings"
	"time"
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
	for k, event := range events {
		l.DataChan <- event
		delete(events, k)
	}

	if l.OpenSysLog {
		log.SetFlags(log.Ldate | log.Lmicroseconds)
		if level == SLevel_Error {
			log.Printf("[%s][%s]: %s\n%s", l.session, level, content, getCaller(3))
		} else {
			log.Printf("[%s][%s]: %s", l.session, level, content)
		}
	}
}

func (l *Logger) getEvents(level string, content string) (events map[string]LoggerEvent) {
	events = make(map[string]LoggerEvent)
	currentLevel := levelMap[level]
	if currentLevel <= levelMap[l.Level] && currentLevel > ILevel_OFF && currentLevel < ILevel_ALL {
		event := LoggerEvent{Level: level, Name: l.Name, Now: time.Now(), Content: content,
			Path: l.Config.Appender.Path, Session: l.session, Caller: getCaller(4)}
		events[level] = event
	}

	/*for k, v := range levelMap {
		if v <= level && v <= logLevel && v > ILevel_OFF {
			event := &LoggerEvent{Level: k, Name: l.Name, Now: time.Now(), Content: content,
				Path: l.Config.Appender.Path}
			events[k] = event
		}
	}*/
	return
}
