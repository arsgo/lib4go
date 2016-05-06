package script

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yuin/gopher-lua"
)

func getHttpModule() (m map[string]lua.LGFunction) {
	return map[string]lua.LGFunction{
		"get":  httpGet,
		"post": httpPost,
	}
}
func httpGet(l *lua.LState) int {
	if l.GetTop() != 1 {
		l.Push(lua.LNil)
		l.Push(lua.LString("函数输入参数错误:http.get"))
	}
	url := l.CheckString(1)
	resp, err := http.Get(url)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}
	l.Push(lua.LString(string(body)))
	l.Push(lua.LNil)
	return 2
}

func httpPost(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.Push(lua.LNil)
		l.Push(lua.LString("函数输入参数错误:http.post"))
	}
	url := l.CheckString(1)
	params := l.CheckString(2)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(params))
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}
	l.Push(lua.LString(string(body)))
	l.Push(lua.LNil)
	return 2
}
