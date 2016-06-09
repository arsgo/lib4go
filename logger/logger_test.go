package logger

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func Testlogger(t *testing.T) {
	lg, err := Get("httpApi", true)
	if err != nil {
		t.Error(err)
	}
	events := lg.getEvents(SLevel_Fatal, "Fatal")
	if len(events) != 3 {
		t.Error("getEvents 返回个数有误", len(events))
	}
	info := events[SLevel_Info]
	if info == nil || info.Level != SLevel_Info {
		t.Error("getEvents 返回的info对象有误")
	}
	errLevel := events[SLevel_Error]
	if errLevel == nil || errLevel.Level != SLevel_Error {
		t.Error("getEvents 返回的debug对象有误")
	}
	fatal := events[SLevel_Fatal]
	if fatal == nil || fatal.Level != SLevel_Fatal {
		t.Error("getEvents 返回的warn对象有误")
	}
}
func TestFileAppender(t *testing.T) {
	format := "2006-01-02 15:04:05"
	now, _ := time.Parse(format, "2016-06-09 13:51:50")
	event := &LoggerEvent{}
	event.Content = "content"
	event.Level = SLevel_Debug
	event.Name = "api"
	event.Now = now
	event.Path = "./logs/%level/%name/%date.log"
	path := getAppendPath(event)
	if !strings.HasSuffix(path, `\logs\debug\api\20160609.log`) {
		t.Error("翻译路径有误", path)
	}
}
func TestWirte(t *testing.T) {
	lg, err := Get("api", false)
	if err != nil {
		t.Error(err)
	}
	lg.SetPath("c:\\t\\%level.log")
	count := 100000

	for k := range levelMap {
		go func(l string) {
			for i := 0; i < count; i++ {
				lg.print(l, fmt.Sprintf("%s:%d", l, i))
			}
		}(k)

	}

	time.Sleep(time.Hour)
}
