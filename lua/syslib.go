package lua

import (
	"github.com/colinyl/lib4go/utility"
	l "github.com/yuin/gopher-lua"
)

func syslibLoader(L *l.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]l.LGFunction{
	"md5":        md5,
	"getLocalIP": getLocalIP,
}

func md5(L *l.LState) int {
	input := L.ToString(1)
	L.Push(l.LString(utility.Md5(input)))
	return 1
}
func getLocalIP(L *l.LState) int {
	input := L.ToString(1)
	L.Push(l.LString(utility.GetLocalIP(input)))
	return 1
}
