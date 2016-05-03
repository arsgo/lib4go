package mq

import (
	lp "github.com/colinyl/lib4go/lua"
	"github.com/yuin/gopher-lua"
)

type ConfigHandler interface {
	GetMQConfig(string) (string, error)
}
type MQBinder struct {
	handler ConfigHandler
	pool    *lp.LuaPool
}

func NewMQBinder(c ConfigHandler, pool *lp.LuaPool) *MQBinder {
	return &MQBinder{handler: c, pool: pool}
}
func (c *MQBinder) BindMQService(L *lua.LState) {
	lp.Bind(L, &lp.ScriptBindClass{ClassName: "mq",
		ConstructorName: "new",
		ConstructorFunc: func(L *lua.LState) interface{} {
			config, _ := c.handler.GetMQConfig(L.CheckString(1))
			s, _ := NewMQService(config)
			return s
		}, ObjectMethods: map[string]lp.ScriptBindFunc{
			"close": func(L *lua.LState) (result []string) {
				if L.GetTop() != 1 {
					result = append(result, "input args error")
					return
				}
				ud := L.CheckUserData(1)
				if _, ok := ud.Value.(IMQService); !ok {
					result = append(result, "MQService expected")
					return
				}
				p := ud.Value.(IMQService)
				p.Close()
				return
			},
			"send": func(L *lua.LState) (result []string) {
				if L.GetTop() != 3 {
					result = append(result, "input args error")
					return
				}
				ud := L.CheckUserData(1)
				if _, ok := ud.Value.(IMQService); !ok {
					result = append(result, "MQService expected")
					return
				}
				p := ud.Value.(IMQService)
				err := p.Send(L.CheckString(2), L.CheckString(3))
				if err != nil {
					result = append(result, err.Error())
				}
				return result
			},
		}})
}
