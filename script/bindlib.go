package script

import (
	"fmt"

	"github.com/colinyl/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

func myfunc() lua.LValue {
	return lua.LString("value")
}
func bindLib(l *lua.LState, binder *LuaBinder) {
	l.SetGlobal("md5", New(l, utility.Md5))
	l.SetGlobal("print", New(l, fmt.Println))
	/*if binder.modeules != nil {
		for k, v := range binder.modeules {
			l.PreloadModule(k, func(l *lua.LState) int {
				return NewModule(l, v)
			})
		}
	}*/
	if binder.libs != nil {
		for k, v := range binder.libs {
			l.SetGlobal(k, New(l, v))
		}
	}
	if binder.types != nil {
		for k, v := range binder.types {
			l.SetGlobal(k, NewType(l, v))
		}
	}

}
