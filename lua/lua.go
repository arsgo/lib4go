package lua

import (
	"errors"

	"github.com/colinyl/lib4go/pool"
	l "github.com/yuin/gopher-lua"
)

type luafunc struct {
	name     string
	function l.LGFunction
}

type luaPoolObject struct {
	state *l.LState
}
type luaPoolFactory struct {
	script string
	count  int
	funcs  map[string]l.LGFunction
}

//LuaPool  LUA对象池
type luaPool struct {
	p     *pool.ObjectPool
	funcs map[string]l.LGFunction
}

var _pool *luaPool

//Close close a object
func (p *luaPoolObject) Close() {
	if p.state != nil {
		p.state.Close()
	}
}

func (l *luaPoolObject) Check() bool {
	return true
}

func (l *luaPoolObject) RequestFatal() {

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
func PreLoad(script string, size int) int {
	return _pool.PreLoad(script, size)
}

//Call 执行脚本main函数
func Call(script string, input ...string) ([]string, error) {
	return _pool.Call(script, input...)
}

//NewLuaPool 构建LUA对象池
func NewLuaPool(funcs ...luafunc) *luaPool {
	cfun := make(map[string]l.LGFunction, 0)
	if len(funcs) > 0 {
		for _, v := range funcs {
			cfun[v.name] = v.function
		}
	}
	return &luaPool{p: pool.New(), funcs: cfun}
}

//PreLoad 预加载脚本
func (p *luaPool) PreLoad(script string, size int) int {
	return p.p.Register(script, &luaPoolFactory{script: script, funcs: p.funcs}, size)
}

//Call 执行脚本main函数
func (p *luaPool) Call(script string, input ...string) ([]string, error) {
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
	var buffer []string
	for _, lv := range values {
		buffer = append(buffer, lv.String())
	}
	return buffer, nil
}
