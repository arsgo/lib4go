package lua

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/colinyl/lib4go/pool"
	l "github.com/yuin/gopher-lua"
)

type Luafunc struct {
	Name     string
	Function l.LGFunction
}

type luaPoolObject struct {
	state *l.LState
}
type luaPoolFactory struct {
	script   string
	count    int
	funcs    map[string]l.LGFunction
	userData []func(*l.LState)
	modules  []Luafunc
}

//LuaPool  LUA对象池
type LuaPool struct {
	p        *pool.ObjectPool
	funcs    map[string]l.LGFunction
	modules  []Luafunc
	userData []func(*l.LState)
}

func (p *LuaPool) AddUserData(f func(*l.LState)) {
	p.userData = append(p.userData, f)
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

//Create create object
func (f *luaPoolFactory) Create() (pool.Object, error) {
	f.count++
	o := &luaPoolObject{}
	o.state = l.NewState()
	o.state.PreloadModule("sys", syslibLoader)
	if f.funcs != nil {
		for k, f := range f.funcs {			
			o.state.PreloadModule(k, f)
		}
	}
	for _, v := range f.userData {
		v(o.state)
	}
	er := o.state.DoFile(f.script)
	if er != nil {
		return nil, er
	}
	return o, nil
}
func (f *luaPoolFactory) registerFunc(name string, fun l.LGFunction, obj *luaPoolObject) {
	obj.state.SetGlobal(name, obj.state.NewFunction(fun))
}

func init() {
	_pool = NewLuaPool()
}

//PreLoad 预加载脚本
func PreLoad(script string, size int) error {
	return _pool.PreLoad(script, size)
}

//Call 执行脚本main函数
func Call(script string, input ...string) ([]string, error) {
	return _pool.Call(script, input...)
}

//NewLuaPool 构建LUA对象池
func NewLuaPool(funcs ...Luafunc) *LuaPool {
	cfun := make(map[string]l.LGFunction, 0)
	if len(funcs) > 0 {
		for _, v := range funcs {
			cfun[v.Name] = v.Function
		}
	}
	return &LuaPool{p: pool.New(), funcs: cfun}
}

//PreLoad 预加载脚本
func (p *LuaPool) PreLoad(script string, size int) error {
	if !exist(script) {
		return errors.New(fmt.Sprintf("not find script :%s", script))
	}
	p.p.Register(script, &luaPoolFactory{script: script, funcs: p.funcs, userData: p.userData,
		modules: p.modules}, size)
	return nil
}

//Call 执行脚本main函数
func (p *LuaPool) Call(script string, input ...string) (result []string, er error) {
	result = []string{}
	if strings.EqualFold(script, "") {
		return result, errors.New(fmt.Sprintf("script(%s) is nil", script))
	}

	if !p.p.Exists(script) {
		er = p.PreLoad(script, 1)
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
	co := L.NewThread() /* create a new thread */
	main := L.GetGlobal("main")
	if main == l.LNil {
		panic(errors.New("cant find main func"))
	}
	fn := main.(*l.LFunction) /* get function from lua */
	var inputs []l.LValue
	for _, v := range input {
		inputs = append(inputs, l.LString(v))
	}
	st, err, values := L.Resume(co, fn, inputs[0:len(input)]...)
	if st == l.ResumeError {
		return nil, err
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
