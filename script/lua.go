package script

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
	"github.com/colinyl/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

type luaPoolObject struct {
	state *lua.LState
}
type luaPoolFactory struct {
	script  string
	binders *LuaBinder
}
type LuaBinder struct {
	packages []string
	libs     map[string]interface{}
	types    map[string]interface{}
	modeules map[string]map[string]lua.LGFunction
}

//LuaPool  LUA对象池
type LuaPool struct {
	p          *pool.ObjectPool
	binders    *LuaBinder
	minSize    int
	maxSize    int
	backStatus chan *lua.LState
	close      chan int
}

var _pool *LuaPool

//Close close a object
func (p *luaPoolObject) Close() {
	if p.state != nil {
		p.state.Close()
	}
}

func (l *luaPoolObject) Check() bool {
	return true
}

func (l *luaPoolObject) Fatal() {

}
func (f *luaPoolFactory) Close() {

}

//Create create object
func (f *luaPoolFactory) Create() (p pool.Object, err error) {
	o := &luaPoolObject{}
	o.state = lua.NewState()
	p = o
	err = bindLib(o.state, f.binders)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.state.DoFile(f.script)
	if err != nil {
		fmt.Println(err)
	}
	return
}

//NewLuaPool 构建LUA对象池
func NewLuaPool() *LuaPool {
	pool := &LuaPool{p: pool.New(), binders: &LuaBinder{}, minSize: 1, maxSize: 10}
	pool.backStatus = make(chan *lua.LState, 100)
	pool.close = make(chan int, 1)
	//	go pool.recycle()
	return pool
}
func (p *LuaPool) recycle() {
	defer p.recycle()
BRK:
	for {
		select {
		case l := <-p.backStatus:
			{
				defer p.recycle()
				//collectgarbage(l)
				l.Close()
				l = nil

			}
		case <-p.close:
			{
				break BRK
			}
		}
	}
}

func (p *LuaPool) SetPackages(paths ...string) {
	if len(paths) == 0 {
		return
	}
	p.binders.packages = append(p.binders.packages, paths...)
}
func (p *LuaPool) RegisterLibs(libs map[string]interface{}) error {
	p.binders.libs = libs
	return nil
}
func (p *LuaPool) RegisterTypes(types map[string]interface{}) error {
	p.binders.types = types
	return nil
}
func (p *LuaPool) RegisterModules(modules map[string]map[string]lua.LGFunction) error {
	p.binders.modeules = modules
	return nil
}

func (p *LuaPool) Close() {
	p.close <- 1
}

//SetPoolSize 设置连接池大小
func (p *LuaPool) SetPoolSize(minSize int, maxSize int) {
	p.minSize = minSize
	p.maxSize = maxSize
}
func (p *LuaPool) getDefSize(m int, def int) int {
	if m <= 0 {
		return def
	}
	return m
}

//PreLoad 预加载脚本
func (p *LuaPool) PreLoad(script string, minSize int, maxSize int) error {
	if !exist(script) {
		return fmt.Errorf("not find script :%s", script)
	}

	p.p.Register(script, &luaPoolFactory{script: script, binders: p.binders}, minSize, maxSize)
	return nil
}
func (p *LuaPool) getScriptLoggerName(name string) string {
	script := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSuffix(strings.ToLower(name), ".lua"), ".luac"), ".")
	rname := strings.Trim(strings.Replace(strings.Replace(script, "/", "-", -1), "\\", "-", -1), "-")

	if index := strings.Index(rname, "script"); index > -1 {
		return strings.Replace(rname[index:], "scripts-", "script/", -1)
	}
	return strings.Replace(rname, "scripts-", "script/", -1)
}

func (p *LuaPool) Call(script string, session string, input string, body string) (result []string, outparams map[string]string, err error) {
	log, err := logger.NewSession(p.getScriptLoggerName(script), session)
	defer luaRecover(log)
	if err != nil {
		return
	}
	result, outparams, err = p.call(script, session, input, body, log)
	if err != nil {
		log.Error(script, err)
	}
	return
}

//GetSnap 获取LUA引擎的快照信息
func (p *LuaPool) GetSnap() pool.ObjectPoolSnap {
	return p.p.GetSnap()
}
func luaRecover(log logger.ILogger) {
	if r := recover(); r != nil {
		log.Fatal(r, string(debug.Stack()))
	}
}

//Call 执行脚本main函数
func (p *LuaPool) call(script string, session string, input string, body string, log logger.ILogger) (result []string, outparams map[string]string, er error) {
	defer luaRecover(log)
	result = []string{}
	if strings.EqualFold(script, "") {
		er = fmt.Errorf("script(%s) is nil", script)
		return
	}
	if !p.p.Exists(script) {
		er = p.PreLoad(script, p.minSize, p.maxSize)
		if er != nil {
			return
		}
	}
	o, er := p.p.Get(script)
	if er != nil {
		return
	}
	co := o.(*luaPoolObject).state
	//log.Info("init top：", co.GetTop())
	dynamicBind(co, map[string]interface{}{
		"print":  log.Info,
		"printf": log.Infof,
		"error":  log.Error,
		"errorf": log.Errorf,
		"fatal":  log.Fatal,
		"fatalf": log.Fatalf,
		"debug":  log.Debug,
		"debugf": log.Debugf,
	})
	defer p.p.Recycle(script, o)
	co.SetGlobal("__session", lua.LString(session))
	main := co.GetGlobal("main")
	if main == lua.LNil {
		er = errors.New("cant find main func")
		return
	}
	outparams = getResponse(co)
	inputArgs := json2LuaTable(co, input, log)
	values, er := callMain(co, inputArgs, lua.LString(body), log)
	for _, lv := range values {
		if strings.EqualFold(lv.Type().String(), "table") {
			result = append(result, luaTable2Json(co, lv, log))
		} else {
			result = append(result, lv.String())
		}
	}
	return
}

//Call 执行脚本main函数
func (p *LuaPool) call2(script string, session string, input string, body string, log logger.ILogger) (result []string, outparams map[string]string, er error) {
	defer luaRecover(log)
	result = []string{}
	if strings.EqualFold(script, "") {
		er = fmt.Errorf("script(%s) is nil", script)
		return
	}
	if !p.p.Exists(script) {
		er = p.PreLoad(script, p.minSize, p.maxSize)
		if er != nil {
			return
		}
	}
	o, er := p.p.Get(script)
	if er != nil {
		return
	}
	L := o.(*luaPoolObject).state
	dynamicBind(L, map[string]interface{}{
		"print":  log.Info,
		"printf": log.Infof,
		"error":  log.Error,
		"errorf": log.Errorf,
		"fatal":  log.Fatal,
		"fatalf": log.Fatalf,
		"debug":  log.Debug,
		"debugf": log.Debugf,
	})
	co := L.NewThread()
	defer p.p.Recycle(script, o)
	defer func(la *lua.LState) {
		p.backStatus <- la
	}(co)
	co.SetGlobal("__session", lua.LString(session))
	main := co.GetGlobal("main")
	if main == lua.LNil {
		er = errors.New("cant find main func")
		return
	}
	outparams = getResponse(co)
	fn := main.(*lua.LFunction)
	inputArgs := json2LuaTable(co, input, log)
	st, err, values := L.Resume(co, fn, inputArgs, lua.LString(utility.Escape(body)))
	if st == lua.ResumeError {
		er = fmt.Errorf("script  error:%s", err)
		return
	}
	for _, lv := range values {
		if strings.EqualFold(lv.Type().String(), "table") {
			result = append(result, luaTable2Json(co, lv, log))
		} else {
			result = append(result, lv.String())
		}
	}
	return
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func collectgarbage(L *lua.LState) error {
	//return L.DoString(`collectgarbage("collect")`)
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

func callMain(ls *lua.LState, inputValue lua.LValue, others lua.LValue, log logger.ILogger) (rt []lua.LValue, er error) {
	defer luaRecover(log)
	ls.Pop(ls.GetTop())
	block := lua.P{
		Fn:      ls.GetGlobal("main"),
		NRet:    2,
		Protect: true,
	}
	er = ls.CallByParam(block, inputValue, others)
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

func json2LuaTable(L *lua.LState, json string, log logger.ILogger) (inputValue lua.LValue) {
	defer luaRecover(log)
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
