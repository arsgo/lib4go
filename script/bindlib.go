package script

import (
	"fmt"
	"strings"

	"github.com/arsgo/lib4go/security/md5"
	"github.com/yuin/gopher-lua"
)

func addPackages(l *lua.LState, paths ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("addPackages:", r)
		}
	}()
	for _, v := range paths {
		pk := `local p = [[` + strings.Replace(v, "//", "/", -1) + `]]
local m_package_path = package.path
package.path = string.format('%s;%s/?.lua;%s/?.luac;%s/?.dll',
	m_package_path, p,p,p)`

		err = l.DoString(pk)
		if err != nil {
			return err
		}
	}

	return
}
func dynamicBind(l *lua.LState, binder map[string]interface{}) {
	for i, v := range binder {
		l.SetGlobal(i, New(l, v))
	}
}

func bindLib(l *lua.LState, binder *LuaBinder) (err error) {
	l.SetGlobal("sys_md5", New(l, md5.Encrypt))
	l.SetGlobal("print", New(l, fmt.Println))

	if binder.packages != nil && len(binder.packages) > 0 {
		err = addPackages(l, binder.packages...)
		if err != nil {
			return
		}
	}
	if binder.modeules != nil {
		for k, v := range binder.modeules {
			l.PreloadModule(k, NewLuaModule(v).Loader)
		}
	}
	if binder.libs != nil {
		for k, v := range binder.libs {
			l.SetGlobal(k, New(l, v))
		}
	}
	if binder.types != nil {
		for _, v := range binder.types {
			mt := l.NewTypeMetatable(v.Name)
			l.SetGlobal(v.Name, mt)
			for i, ff := range v.NewFunc {
				l.SetField(mt, i, l.NewFunction(ff))
			}
			l.SetField(mt, "__index", l.SetFuncs(l.NewTable(), v.Methods))
		}
	}
	if binder.global != nil {
		for i, v := range binder.global {
			l.SetGlobal(i, l.NewFunction(v))
		}
	}
	return
}
