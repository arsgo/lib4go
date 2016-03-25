package lua

import (
	"strings"
	"testing"
    "strconv"
)

func Test_lua(t *testing.T) {
	luaPool := NewLuaPool()
	if luaPool.PreLoad("./t1.lua", 2)!=2{
        t.Error("luapool init error")
    }
    if luaPool.PreLoad("./t2.lua", 2)!=2{
        t.Error("luapool init error")
    }
	values, err := luaPool.Call("./t1.lua")
	if err != nil {
		t.Error(err.Error())
	}
	for i, v := range values {
		if !strings.EqualFold(strconv.Itoa(i+1), string(v)) {
			t.Errorf("return values is error [%s]:[%s]", strconv.Itoa(i+1), string(v))
		}
	}    
	values, err = luaPool.Call("./t2.lua", "123456")
	if err != nil {
		t.Error(err.Error())
	}
    if len(values)!=1{
        t.Error("return values len error")
    }
	for _, v := range values {
		if !strings.EqualFold("e10adc3949ba59abbe56e057f20f883e", v) {
			t.Errorf("return values is error %s", v)
		}
	}
    
    if PreLoad("./t1.lua", 2)!=2{
        t.Error("luapool init error")
    }
    if PreLoad("./t2.lua", 2)!=2{
        t.Error("luapool init error")
    }
	values, err = Call("./t1.lua")
	if err != nil {
		t.Error(err.Error())
	}
	for i, v := range values {
		if !strings.EqualFold(strconv.Itoa(i+1), string(v)) {
			t.Errorf("return values is error [%s]:[%s]", strconv.Itoa(i+1), string(v))
		}
	}    
	values, err = Call("./t2.lua", "123456")
	if err != nil {
		t.Error(err.Error())
	}
    if len(values)!=1{
        t.Error("return values len error")
    }
	for _, v := range values {
		if !strings.EqualFold("e10adc3949ba59abbe56e057f20f883e", v) {
			t.Errorf("return values is error %s", v)
		}
	}
    
    
}
