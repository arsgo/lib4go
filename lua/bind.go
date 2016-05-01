package lua

import l "github.com/yuin/gopher-lua"

type ScriptBindFunc func(*l.LState) []string
type ScriptBinderFuncs map[string]ScriptBindFunc

type ScriptBindClass struct {
	ClassName       string
	ConstructorName string
	ConstructorFunc func(*l.LState) interface{}
	Funcs           ScriptBinderFuncs
	ObjectMethods   ScriptBinderFuncs
}

func Bind(L *l.LState, pk *ScriptBindClass) {
	mt := L.NewTypeMetatable(pk.ClassName)
	L.SetGlobal(pk.ClassName, mt)
	L.SetField(mt, pk.ConstructorName, L.NewFunction(func(ls *l.LState) int {
		ud := ls.NewUserData()
		ud.Value = pk.ConstructorFunc(ls)
		ls.SetMetatable(ud, ls.GetTypeMetatable(pk.ClassName))
		ls.Push(ud)
		return 1
	}))

	for cName, cFunc := range pk.Funcs {
		L.SetField(mt, cName, L.NewFunction(getFunc(L, cFunc)))
	}
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), getFuncMap(L, pk.ObjectMethods)))
}
func getFunc(L *l.LState, fun ScriptBindFunc) l.LGFunction {
	return func(ls *l.LState) int {
		results := fun(ls)
		for _, v := range results {
			ls.Push(l.LString(v))
		}
		return len(results)
	}
}
func getFuncMap(L *l.LState, funs ScriptBinderFuncs) (rfun map[string]l.LGFunction) {
	rfun = make(map[string]l.LGFunction)
	for i, v := range funs {
		rfun[i] = getFunc(L, v)
	}
	return
}
