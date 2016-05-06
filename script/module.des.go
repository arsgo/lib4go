package script

import (
	"fmt"

	"github.com/colinyl/lib4go/des"
	"github.com/yuin/gopher-lua"
)

func getDesModule() (m map[string]lua.LGFunction) {
	return map[string]lua.LGFunction{
		"encrypt": desEncrypt,
		"decrypt": desDecrypt,
	}
}
func desDecrypt(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.Push(lua.LNil)
		l.Push(lua.LString("函数输入参数错误:desc.encrypt"))
	}
	origin := l.CheckString(1)
	key := l.CheckString(2)
	sec, err := des.Decrypt(origin, key)
	l.Push(lua.LString(sec))
	l.Push(lua.LString(fmt.Sprintf("%v", err)))
	return 2
}

func desEncrypt(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.Push(lua.LNil)
		l.Push(lua.LString("函数输入参数错误:desc.decrypt"))
	}
	origin := l.CheckString(1)
	key := l.CheckString(2)
	sec, err := des.Encrypt(origin, key)
	l.Push(lua.LString(sec))
	l.Push(lua.LString(fmt.Sprintf("%v", err)))
	return 2
}
