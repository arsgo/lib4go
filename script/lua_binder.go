package script

import lua "github.com/yuin/gopher-lua"

type lualoader struct {
	script string
	min    int
	max    int
}

type LuaTypesBinder struct {
	Name    string
	NewFunc map[string]lua.LGFunction
	Methods map[string]lua.LGFunction
}

type LuaBinder struct {
	packages []string
	libs     map[string]interface{}
	types    []LuaTypesBinder
	global   map[string]lua.LGFunction
	modeules map[string]map[string]lua.LGFunction
	preload  []lualoader
}

//SetPackages 设置LUA引擎的基础package
func (p *LuaBinder) SetPackages(paths ...string) {
	if len(paths) == 0 {
		return
	}
	p.packages = append(p.packages, paths...)
}

//RegisterScript 注册脚本
func (p *LuaBinder) RegisterScript(script string, min int, max int) {
	p.preload = append(p.preload, lualoader{script: script, min: min, max: max})
}

//RegisterGlobal 注册公共函数
func (p *LuaBinder) RegisterGlobal(values map[string]lua.LGFunction) error {
	p.global = values
	return nil
}

//RegisterLibs 注册引用的lib路径
func (p *LuaBinder) RegisterLibs(libs map[string]interface{}) error {
	p.libs = libs
	return nil
}

//RegisterTypes 注册GO类型供lua调用
func (p *LuaBinder) RegisterTypes(types ...LuaTypesBinder) error {
	p.types = types
	return nil
}

func (p *LuaBinder) RegisterModules(modules map[string]map[string]lua.LGFunction) error {
	p.modeules = modules
	return nil
}
