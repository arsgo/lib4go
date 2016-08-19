package script

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/pool"
	lua "github.com/yuin/gopher-lua"
)

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
func (p *luavm) Call(script string, session string, input string, body string) (result []string, outparams map[string]string, err error) {
	log, err := logger.NewSession(getScriptLoggerName(script), session)
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

//Call 执行脚本main函数
func (p *luavm) call(script string, session string, input string, body string, log logger.ILogger) (result []string, outparams map[string]string, er error) {
	defer luaRecover(log)
	result = []string{}
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
	defer p.p.Recycle(script, o)
	co.SetGlobal("__session", lua.LString(session))
	co.SetGlobal("__logger_name", lua.LString(log.GetName()))
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
func (p *luavm) Close() {
	p.p.Close()
}
