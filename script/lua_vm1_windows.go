package script
/*
import (
	"fmt"

	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/pool"
	lua "github.com/yuin/gopher-lua"
)

type InputArgs struct {
	Script   string
	Session  string
	Input    string
	Body     string
	TaskName string
	TaskType string
	recvChan chan dataMessage
	log      logger.ILogger
}
type luaStateVM struct {
	state    *lua.LState
	msgChan  chan lua.LValue
	quitChan chan lua.LValue
	count    int32
}
type luaResult struct {
	vm      *luaStateVM
	session string
	result  []string
}
type dataMessage struct {
	result []string
}

type luavm struct {
	states        []*luaStateVM
	Binder        *LuaBinder
	sysMsgChan    chan InputArgs
	sysResultChan chan lua.LValue
	dataMap       map[string]chan dataMessage
	log           logger.ILogger
	minSize       int
	maxSize       int
}

func newLuaVM(binder *LuaBinder, minSize int, maxSize int) *luavm {
	vm := &luavm{Binder: binder, minSize: minSize, maxSize: maxSize}
	vm.log, _ = logger.Get("app.server")
	vm.sysMsgChan = make(chan InputArgs, 1000)
	vm.sysResultChan = make(chan lua.LValue, 1000)
	vm.dataMap = make(map[string]chan dataMessage)
	vm.states = make([]*luaStateVM, 0, 8)
	go vm.dispatch()
	return vm
}
func (p *luavm) dispatch() {
	for {
		select {
		case lv := <-p.sysResultChan: //系统通道接口消息，并发送运行时通道
			r := lv.(*lua.LTable)
			session := lua.LVAsString(r.RawGetString("session"))
			fmt.Println("recv.data:", session)
			if data, ok := p.dataMap[session]; ok {
				data1 := lua.LVAsString(r.RawGetString("data1"))
				data2 := lua.LVAsString(r.RawGetString("data2"))
				data <- dataMessage{result: []string{data1, data2}}
			}
			ndata := <-p.sysMsgChan //从系统通道获取任务，并分配给运行通道
			ndata.log.Info("abc..1", ndata.Session)
			vm := r.RawGetString("vm").(*lua.LUserData).Value.(*luaStateVM)

			vm.state.SetGlobal("__session__", lua.LString(ndata.Session))
			vm.state.SetGlobal("__logger_name__", lua.LString(ndata.log.GetName()))
			vm.state.SetGlobal("__task_name__", lua.LString(ndata.TaskName))
			vm.state.SetGlobal("__task_type__", lua.LString(ndata.TaskType))

			input := vm.state.NewTable()                                  //输入参数，传入请求串，回传通道
			input.RawSetString("callback", lua.LChannel(p.sysResultChan)) //用于接口回传参数
			input.RawSetString("input", json2LuaTable(vm.state, ndata.Input, ndata.log))
			input.RawSetString("body", lua.LString(ndata.Body))
			input.RawSetString("vm", New(vm.state, vm))
			input.RawSetString("session", lua.LString(ndata.Session))
			ndata.log.Info("abc..2", ndata.Session, ndata.Input)
			vm.msgChan <- input
		}
	}
}

//PreLoad 预加载脚本
func (p *luavm) PreLoad(script string, minSize int, maxSize int) (err error) {
	vm := &luaStateVM{}
	vm.state = lua.NewState()
	vm.quitChan = make(chan lua.LValue, 1)
	vm.msgChan = make(chan lua.LValue, 1000)
	err = bindLib(vm.state, p.Binder)
	if err != nil {
		vm.state.Close()
		return
	}
	vm.state.SetGlobal("__input_args_chan__", lua.LChannel(vm.msgChan))
	vm.state.SetGlobal("__quit_chan_", lua.LChannel(vm.quitChan))
	err = vm.state.DoFile(script)
	if err != nil {
		vm.state.Close()
	}
	go func(vm *luaStateVM) {
		err := vm.state.DoString(`	 
    local exit = false
    while not exit do
      channel.select(
        {"|<-", __input_args_chan__, function(ok, args)			
          if not ok then
            exit = true
          else       	
          local result={}
		  local r1,r2=main(args.input,args.body) 
          result["data1"]= r1
		  result["data2"]= r2
		  result["vm"]=args.vm
		  result["session"]=args.session		
          args.callback:send(result)
		  print("result:",r1,r2,args.session)
          end
        end},
        {"|<-", __quit_chan_, function(ok, v)        
            exit = true
        end}
      )
    end
    `)
		fmt.Println("err:", err)
	}(vm)
	p.states = append(p.states, vm)
	tb := vm.state.NewTable()
	tb.RawSetString("session", lua.LNil)
	tb.RawSetString("vm", New(vm.state, vm))
	p.sysResultChan <- tb
	return
}

//Call 调用脚本引擎并执行main函数
func (p *luavm) Call(input InputArgs) (result []string, outparams map[string]string, err error) {
	log, err := logger.NewSession(getScriptLoggerName(input.Script), input.Session)
	log.Infof("call:%+v\n", input)
	input.recvChan = make(chan dataMessage, 1)
	p.dataMap[input.Session] = input.recvChan
	input.log = log
	p.sysMsgChan <- input
	value := <-input.recvChan
	result = value.result
	return
}
func (p *luavm) Close() {

}
func (p *luavm) GetSnap() pool.ObjectPoolSnap {
	return pool.ObjectPoolSnap{}
}
*/