package pool

import (
	"fmt"
	"testing"
)

func Test_pool(t *testing.T) {
	var groupName string = "connect_server"
	pools := New()
	_, e := pools.Get(groupName)
	if e == nil {
		t.Error("get error")
	}

	var size int = 10
	var poolSize int = len(pools.pools)
	pools.Register(groupName, &NetworkFactory{}, size)
	if len(pools.pools) != poolSize+1 {
		t.Error("pools register errors")
	}

	if pools.pools[groupName].Size != size {
		t.Error("pools register size errors")
	}

	o, err := pools.Get(groupName)
	if err != nil {
		t.Error(err)
	}

	if pools.pools[groupName].list.Len() != size-1 {
		t.Error(fmt.Sprintf("pool size get error:e:%d,a:%d", (size - 1), pools.pools[groupName].list.Len()))
	}

	pools.Recycle(groupName, o)
	if pools.pools[groupName].list.Len() != size {
		t.Error(fmt.Sprintf("pool size recycle error:e:%d,a:%d", size, pools.pools[groupName].list.Len()))
	}
   
	pools.Close(groupName)
	_, err= pools.Get(groupName)
	if err == nil {
		t.Error("get error")
	}
    
}
