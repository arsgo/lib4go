package script

import (
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

var luaPool *LuaPool
var min int
var max int

func init() {
	min = 100
	max = 1000
	luaPool = NewLuaPool()
}
func TestInit(t *testing.T) {
	if luaPool.PreLoad("./t1.lua", min, max) != nil {
		t.Error("luapool init error")
	}
	if luaPool.PreLoad("./t2.lua", min, max) != nil {
		t.Error("luapool init error")
	}
}
func TestBenchCall(t *testing.T) {
	time.Sleep(time.Second * 2)
	ch := make(chan int, max)
	close := make(chan int, 1)
	var index int32
	var concurrent int32
	concurrent = 100000
	groupName := "./t2.lua"

	for i := 0; i < min; i++ {
		ch <- i
		go func() {
			for {
				if atomic.LoadInt32(&index) >= concurrent {
					close <- 1
					break
				}
				<-ch
				values, err := luaPool.Call(groupName, "123456")
				if err != nil {
					t.Error(err.Error())
				} else {
					if len(values) != 1 {
						t.Error("return values len error")
					}
				}
				atomic.AddInt32(&index, 1)
				ch <- 1
			}

		}()
	}
	<-close

}

func TestLua(t *testing.T) {

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
	if len(values) != 1 {
		t.Error("return values len error")
	}
	for _, v := range values {
		if !strings.EqualFold("e10adc3949ba59abbe56e057f20f883e", v) {
			t.Errorf("return values is error %s", v)
		}
	}

}
