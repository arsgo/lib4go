package script

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/utility"
	lua "github.com/yuin/gopher-lua"
)

func luaRecover(log logger.ILogger) {
	if r := recover(); r != nil {
		log.Fatal(r, string(debug.Stack()))
	}
}
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func collectgarbage(L *lua.LState) error {
	return L.DoString(`collectgarbage("setpause"，10)`)
}
func getResponse(L *lua.LState) (r map[string]string) {
	fields := map[string]string{
		"content_type": "Content-Type",
		"charset":      "Charset",
		"original":     "_original",
	}
	r = make(map[string]string)
	response := L.GetGlobal("response")
	if response == lua.LNil {
		return
	}
	for i, v := range fields {
		fied := L.GetField(response, i)
		if fied == lua.LNil {
			continue
		}
		r[v] = fied.String()
	}
	return
}
func getScriptLoggerName(name string) string {
	script := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSuffix(strings.ToLower(name), ".lua"), ".luac"), ".")
	rname := strings.Trim(strings.Replace(strings.Replace(script, "/", "-", -1), "\\", "-", -1), "-")

	if index := strings.Index(rname, "script"); index > -1 {
		return strings.Replace(rname[index:], "scripts-", "script/", -1)
	}
	return strings.Replace(rname, "scripts-", "script/", -1)
}
func callMainFunc(ls *lua.LState, args ...lua.LValue) (err error) {
	defer base.RunTime("call script  main", time.Now())
	block := lua.P{
		Fn:      ls.GetGlobal("main"),
		NRet:    2,
		Protect: true,
	}
	err = ls.CallByParam(block, args...)
	return err
}

func callMain(ls *lua.LState, inputValue lua.LValue, others lua.LValue, log logger.ILogger) (rt []lua.LValue, er error) {
	defer luaRecover(log)
	ls.Pop(ls.GetTop())
	er = callMainFunc(ls, inputValue, others)
	if er != nil {
		return
	}
	defer ls.Pop(ls.GetTop())
	rt = make([]lua.LValue, 0, ls.GetTop())
	value1 := ls.Get(1)
	if value1.String() == "nil" {
		return
	}
	rt = append(rt, value1)
	if value1.String() != "302" {
		return
	}
	value2 := ls.Get(2)
	if value2.String() == "nil" {
		return
	}
	rt = append(rt, value2)
	return
}

func luaTable2Json(L *lua.LState, inputValue lua.LValue, log logger.ILogger) (json string) {
	defer luaRecover(log)
	defer base.RunTime("call script  table2json", time.Now())
	L.Pop(L.GetTop())
	xjson := L.GetGlobal("xjson")
	if xjson.String() == "nil" {
		fmt.Println("not find xjson")
		json = inputValue.String()
		return
	}
	encode := L.GetField(xjson, "encode")
	if encode == nil {
		fmt.Println("not find xjson.encode")
		json = inputValue.String()
		return
	}
	block := lua.P{
		Fn:      encode,
		NRet:    1,
		Protect: true,
	}
	er := L.CallByParam(block, inputValue)
	if er != nil {
		fmt.Println(er)
		json = inputValue.String()
	} else {
		json = L.Get(-1).String()
	}
	L.Pop(L.GetTop())
	return
}
func json2LuaTable(L *lua.LState, input string, log logger.ILogger) (inputValue lua.LValue) {
	defer luaRecover(log)
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		inputValue = lua.LString(input)
		return
	}
	tb := L.NewTable()
	for k, v := range data {
		tb.RawSetString(k, json2LuaTableValue(L, v, log))
	}
	return tb
}
func json2LuaTableValue(L *lua.LState, value interface{}, log logger.ILogger) (inputValue lua.LValue) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Slice:
		nvalue := value.([]interface{})
		tb := L.NewTable()
		for k, v := range nvalue {
			tb.RawSetInt(k, json2LuaTableValue(L, v, log))
		}
		return tb
	case reflect.Map:
		nvalue := value.(map[string]interface{})
		tb := L.NewTable()
		for k, v := range nvalue {
			tb.RawSetString(k, json2LuaTableValue(L, v, log))
		}
		return tb
	default:
		inputValue = New(L, value)
	}
	return
}

func json2LuaTable2(L *lua.LState, json string, log logger.ILogger) (inputValue lua.LValue) {
	defer luaRecover(log)
	defer base.RunTime("call script  json2table", time.Now())
	L.Pop(L.GetTop())
	xjson := L.GetGlobal("xjson")
	if xjson.String() == "nil" {
		inputValue = lua.LString(json)
		return
	}
	decode := L.GetField(xjson, "decode")
	if decode == nil {
		inputValue = lua.LString(utility.Escape(json))
		return
	}
	block := lua.P{
		Fn:      decode,
		NRet:    1,
		Protect: true,
	}
	er := L.CallByParam(block, lua.LString(utility.Escape(json)))
	if er != nil {
		inputValue = lua.LString(json)
	} else {
		inputValue = L.Get(-1)
	}
	L.Pop(L.GetTop())
	return
}
