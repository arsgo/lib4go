package pool

import (
	"fmt"
	"testing"
)

type NetworkClient struct {
	index int
}

func (n *NetworkClient) Close() {
	fmt.Printf("index:%d", n.index)
}

type NetworkFactory struct {
	index int
}

func (n *NetworkFactory) Create() Object {
	n.index++
	return &NetworkClient{index: n.index}
}

func Test_poolset(t *testing.T) {
	var count int = 10
	set := newPoolSet(count, count, &NetworkFactory{})
	if set.list.Len() != count {
		t.Error(fmt.Sprintf("初始化失败:%d", set.list.Len()))
	}

	if set.Size != set.list.Len() {
		t.Error("pool set size error")
	}

	for i := 1; i <= count; i++ {
		o, _ := set.get()
		v := o.(*NetworkClient)
		if set.list.Len() != count-1 {
			t.Error("长度有误")
		}
		set.add(o)
		if set.list.Len() != count {
			t.Error("长度有误")
		}
		if v.index != i {
			t.Error(fmt.Sprintf("不相待:%d,%d", v.index, i))
		}
	}
}
