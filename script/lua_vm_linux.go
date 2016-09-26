package script
/*
import (
	"errors"
	"strings"
	"sync"

	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/pool"
	lua "github.com/yuin/gopher-lua"
)

type InputArgs struct {
	Script      string
	Content     string
	Session     string
	Input       string
	Body        string
	TaskName    string
	TaskType    string
	HTTPContext interface{}
}

type luavm struct {
	locker   sync.Mutex
	pool     []*lua.LState
	Binder   *LuaBinder
	isClosed bool
	minSize  int
	maxSize  int
}

func newLuaVM(binder *LuaBinder, minSize int, maxSize int) *luavm {
	v := &luavm{Binder: binder, minSize: minSize, maxSize: maxSize, isClosed: false}
	v.pool = make([]*lua.LState, 0, maxSize)
	for i := 0; i < minSize; i++ {
		v.new()
	}
	return v
}

func (p *luavm) GetSnap() pool.ObjectPoolSnap {
	return pool.ObjectPoolSnap{}
}

func (p *luavm) new() (state *lua.LState, err error) {
	state = lua.NewState()
	err = bindLib(state, p.Binder)
	if err != nil {
		state.Close()
		return
	}
	p.pool = append(p.pool, state)
	return
}
func (p *luavm) Get() (state *lua.LState, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	n := len(p.pool)
	if n == 0 {
		state, err = p.new()
		return
	}
	state = p.pool[n-1]
	p.pool = p.pool[0 : n-1]
	return
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
func (p *luavm) Put(L *lua.LState) {
	p.locker.Lock()
	defer p.locker.Unlock()
	if p.isClosed {
		L.Close()
		return
	}
	n := len(p.pool)
	if n < p.maxSize {
		p.pool = append(p.pool, L)
		return
	}
	L.Close()
}

//Call 执行脚本main函数
func (p *luavm) call(input InputArgs, log logger.ILogger) (result []string, outparams map[string]string, er error) {
	defer luaRecover(log)
	result = []string{}

	co, er := p.Get()
	if er != nil {
		return
	}
	defer p.Put(co)
	//er = co.DoFile(input.Script)
	er = co.DoString(input.Content)
	if er != nil {
		return
	}
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
	p.locker.Lock()
	defer p.locker.Unlock()
	p.isClosed = true
	for _, L := range p.pool {
		L.Close()
	}
}
*/