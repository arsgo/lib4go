package script

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/pool"
	lua "github.com/yuin/gopher-lua"
)

type InputArgs struct {
	Script      string
	Session     string
	Input       string
	Body        string
	TaskName    string
	TaskType    string
	HTTPContext interface{}
	MainLogger  logger.ILogger
}

type luavm struct {
	p       *pool.ObjectPool
	Binder  *LuaBinder
	minSize int
	maxSize int
}

func newLuaVM(binder *LuaBinder, minSize int, maxSize int) *luavm {
	return &luavm{p: pool.New(), Binder: binder, minSize: minSize, maxSize: maxSize}
}

func (p *luavm) GetSnap() pool.ObjectPoolSnap {
	return p.p.GetSnap()
}

//PreLoad 预加载脚本
func (p *luavm) PreLoad(script string, minSize int, maxSize int) error {
	if !exist(script) {
		return fmt.Errorf("未找到脚本 :[%s]", script)
	}
	p.p.Register(script, &luaPoolFactory{script: script, binders: p.Binder}, minSize, maxSize)
	return nil
}

//Call 调用脚本引擎并执行main函数
func (p *luavm) Call(input InputArgs) (result []string, outparams map[string]string, err error) {
	log, err := logger.NewSession(getScriptLoggerName(input.Script), input.Session)
	//defer base.RunTime("call script  time", time.Now())
	defer luaRecover(log)
	if err != nil {
		return
	}
	result, outparams, err = p.call(input, log)
	if err != nil {
		log.Error(input.Script, err)
	}
	return
}

//Call 执行脚本main函数
func (p *luavm) call(input InputArgs, log logger.ILogger) (result []string, outparams map[string]string, er error) {
	defer luaRecover(log)
	result = []string{}
	if !p.p.Exists(input.Script) {
		er = p.PreLoad(input.Script, p.minSize, p.maxSize)
		if er != nil {
			return
		}
	}
	o, er := p.p.Get(input.Script)
	if er != nil {
		return
	}
	co := o.(*luaPoolObject).state
	defer p.p.Recycle(input.Script, o)
	co.SetGlobal("__session__", lua.LString(input.Session))
	co.SetGlobal("__logger_name__", lua.LString(log.GetName()))
	co.SetGlobal("__task_name__", lua.LString(input.TaskName))
	co.SetGlobal("__task_type__", lua.LString(input.TaskType))
	co.SetGlobal("__http_context__", New(co, input.HTTPContext))
	co.SetGlobal("__set_cookie__", lua.LNil)
	main := co.GetGlobal("main")
	if main == lua.LNil {
		er = errors.New("cant find main func")
		return
	}

	inputData := json2LuaTable(co, input.Input, log)
	values, er := callMain(co, inputData, lua.LString(input.Body), log)
	for _, lv := range values {
		if strings.EqualFold(lv.Type().String(), "table") {
			result = append(result, luaTable2Json(co, lv, log))
		} else {
			result = append(result, lv.String())
		}
	}
	outparams = getResponse(co)
		return
}
func (p *luavm) Close() {
	p.p.Close()
}
