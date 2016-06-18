package script

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/colinyl/lib4go/pool"
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
	p       *pool.ObjectPool
	binders *LuaBinder
	minSize int
	maxSize int
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
	return &LuaPool{p: pool.New(), binders: &LuaBinder{}, minSize: 1, maxSize: 10}
}
func (p *LuaPool) SetPackages(paths ...string) {
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

//SetPoolSize 设置连接池大小
func (p *LuaPool) SetPoolSize(minSize int, maxSize int) {
	p.minSize = minSize
	p.maxSize = maxSize
}

//PreLoad 预加载脚本
func (p *LuaPool) PreLoad(script string, minSize int, maxSize int) error {
	if !exist(script) {
		return errors.New(fmt.Sprintf("not find script :%s", script))
	}
	p.p.Register(script, &luaPoolFactory{script: script, binders: p.binders}, minSize, maxSize)
	return nil
}

//Call 执行脚本main函数
func (p *LuaPool) Call(script string, input ...string) (result []string, er error) {
	result = []string{}
	if strings.EqualFold(script, "") {
		return result, errors.New(fmt.Sprintf("script(%s) is nil", script))
	}
	if !p.p.Exists(script) {
		er = p.PreLoad(script, p.minSize, p.maxSize)
		if er != nil {
			return
		}
	}
	o, er := p.p.Get(script)
	if er != nil {
		return nil, er
	}
	defer p.p.Recycle(script, o)
	L := o.(*luaPoolObject).state
	co := L.NewThread()
	main := L.GetGlobal("main")
	if main == lua.LNil {
		return nil, errors.New("cant find main func")
	}
	fn := main.(*lua.LFunction)
	var inputs []lua.LValue
	for _, v := range input {
		inputs = append(inputs, lua.LString(v))
	}
	st, err, values := L.Resume(co, fn, inputs[0:len(input)]...)
	co.Close()
	if st == lua.ResumeError {
		return nil, fmt.Errorf("resume error:%s", err)
	}
	for _, lv := range values {
		result = append(result, lv.String())
	}
	return
}
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
