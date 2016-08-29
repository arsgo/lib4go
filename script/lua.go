package script

import (
	"fmt"

	"github.com/arsgo/lib4go/pool"
	"github.com/yuin/gopher-lua"
)

type luaPoolObject struct {
	state *lua.LState
}
type luaPoolFactory struct {
	script  string
	binders *LuaBinder
}

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
		o.state.Close()
		fmt.Println(err)
		return
	}
	err = o.state.DoFile(f.script)
	if err != nil {
		o.state.Close()
		fmt.Println(err)
	}
	return
}
