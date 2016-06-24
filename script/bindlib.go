package script

import (
	"fmt"

	"github.com/colinyl/lib4go/des"
	"github.com/colinyl/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

func addPackages(l *lua.LState, paths ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("addPackages:", r)
		}
	}()
	for _, v := range paths {
		pk := `local p = [[` + v + `]]
local m_package_path = package.path
package.path = string.format('%s;%s/?.lua;%s/?.luac;%s/?.dll',
	m_package_path, p,p,p)`
		//	fmt.Println("pk.path:", pk)

		err = l.DoString(pk)
		if err != nil {
			return err
		}
	}

	return
}
func bindLib(l *lua.LState, binder *LuaBinder) (err error) {
	l.SetGlobal("sys_md5", New(l, utility.Md5))
	l.SetGlobal("print", New(l, fmt.Println))
	l.SetGlobal("des_encrypt", New(l, des.Encrypt))
	l.SetGlobal("des_decrypt", New(l, des.Decrypt))
	//	l.PreloadModule("des", NewLuaModule(getDesModule()).Loader)
	//	l.PreloadModule("http", NewLuaModule(getHttpModule()).Loader)

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
		for k, v := range binder.types {
			l.SetGlobal(k, NewType(l, v))
		}
	}
	return
}
