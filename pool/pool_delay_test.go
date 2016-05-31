package pool

import (
	"fmt"
	"testing"
	"time"
)

type delayClient struct {
}

func (n *delayClient) Close() {
}

type delayFactory struct {
}

func (n *delayFactory) Create() (Object, error) {
	time.Sleep(time.Second)
	return &delayClient{}, nil
}
func (n *delayFactory) Close() {

}
func checkPoolSet(p interface{}, current int32, canUse int32) error {
	cur := p.(*poolSet).current
	use := p.(*poolSet).canUse
	if cur != current || use != canUse {
		return fmt.Errorf("对象数量有误:current:expect-%d,actual-%d, canUse:expect-%d,actual-%d",
			current, cur, canUse, use)
	}
	return nil

}

func TestpoolSetCreate(t *testing.T) {
	groupName := "delay_factory"
	max := 10
	delayPool := New()
	delayPool.Register(groupName, &delayFactory{}, 1, max)

	//------未创建完成时获取对象应失败
	_, err := delayPool.Get(groupName)
	if err == nil {
		t.Error("对象应延迟创建")
	}
	err = checkPoolSet(delayPool.pools.Get(groupName), 0, 0)
	if err != nil {
		t.Error(err)
	}

	//------休息指定时间，应成功获取对象
	time.Sleep(time.Second * 2)
	err = checkPoolSet(delayPool.pools.Get(groupName), 1, 1)
	if err != nil {
		t.Error(err)
	}
	obj, err := delayPool.Get(groupName)
	if err != nil {
		t.Error("对象延迟获取失败，", err)
	}
	err = checkPoolSet(delayPool.pools.Get(groupName), 1, 0)
	if err != nil {
		t.Error(err)
	}

	//------重复获取应返回失败
	_, err = delayPool.Get(groupName)
	if err == nil {
		t.Error("对象创建个数有误")
	}

	//回收对象后，可用数应增加
	delayPool.Recycle(groupName, obj)
	err = checkPoolSet(delayPool.pools.Get(groupName), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_, err = delayPool.Get(groupName)
	if err != nil {
		t.Error("对象回收有误")
	}
	err = checkPoolSet(delayPool.pools.Get(groupName), 1, 0)
	if err != nil {
		t.Error(err)
	}

	//标记为不可用时，可用数应减少
	delayPool.Unusable(groupName, obj)
	err = checkPoolSet(delayPool.pools.Get(groupName), 0, 0)
	if err != nil {
		t.Error(err)
	}
	_, err = delayPool.Get(groupName)
	if err == nil {
		t.Error("对象未正确标记为错误")
	}

	//等待一定时间后可用数应增加
	time.Sleep(time.Second * 2)
	_, err = delayPool.Get(groupName)
	if err != nil {
		t.Error("对象二次创建有误")
	}
	err = checkPoolSet(delayPool.pools.Get(groupName), 1, 0)
	if err != nil {
		t.Error(err)
	}
	delayPool.Recycle(groupName, obj)
	err = checkPoolSet(delayPool.pools.Get(groupName), 1, 1)
	if err != nil {
		t.Error(err)
	}

}
func TestpoolSetBenchInit(t *testing.T) {
	groupName := "delay_factory"
	max := 10
	max32 := int32(max)
	delayPool := New()
	delayPool.Register(groupName, &delayFactory{}, max, max)
	err := checkPoolSet(delayPool.pools.Get(groupName), 0, 0)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * time.Duration(max) * 2)
	err = checkPoolSet(delayPool.pools.Get(groupName), max32, max32)
	if err != nil {
		t.Error(err)
	}

}

func TestpoolSetBenchCreate(t *testing.T) {
	groupName := "delay_factory"
	max := 10
	max32 := int32(max)
	delayPool := New()
	delayPool.Register(groupName, &delayFactory{}, 1, max)
	for i := 0; i < max; i++ {
		_, err := delayPool.Get(groupName)
		if err == nil {
			t.Error("对象批量创建有误")
		}
	}
	err := checkPoolSet(delayPool.pools.Get(groupName), 0, 0)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * time.Duration(max) * 2)
	err = checkPoolSet(delayPool.pools.Get(groupName), max32, max32)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < max; i++ {
		_, err := delayPool.Get(groupName)
		if err != nil {
			t.Error("对象延迟创建失败")
		}
	}
	err = checkPoolSet(delayPool.pools.Get(groupName), max32, 0)
	if err != nil {
		t.Error(err)
	}
	_, err = delayPool.Get(groupName)
	if err == nil {
		t.Error("请求对象超过最大值，应返回失败")
	}
	err = checkPoolSet(delayPool.pools.Get(groupName), max32, 0)
	if err != nil {
		t.Error(err)
	}
}
