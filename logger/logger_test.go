package logger

import (
	"testing"
    "time"
)
func Test_logger(t *testing.T) {
   log,e:=New("httpApi")
   if e!=nil{
       t.Error(e)
   }
   log.Info("hello")
    time.Sleep(time.Second*3)
}