package logger

import (
	"github.com/colinyl/lib4go/concurrent"
)

var logMap *concurrent.ConcurrentMap

func init() {
	logMap = concurrent.NewConcurrentMap()
}

/*
//GetDeubgLogger 获取用于调试的日志
func GetDeubgLogger(session string) ILogger {
	const key string = "flow"
	lg, _ := NewSession(key, session)
	return logMap.GetOrAdd(key, lg).(ILogger)
}
*/
