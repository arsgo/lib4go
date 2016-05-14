package pool

import (
	"errors"
	"fmt"

	"github.com/colinyl/lib4go/concurrent"
)

//ObjectPool 对象缓存池, 缓存池中的对象只添加,不会修改或删除,部分代码对锁进行了优化
type ObjectPool struct {
	pools concurrent.ConcurrentMap
}

//New 创建一个新的对象k
func New() *ObjectPool {
	pools := &ObjectPool{}
	pools.pools = concurrent.NewConcurrentMap()
	return pools
}

//Register 注册指定的对象
func (p *ObjectPool) Register(name string, factory ObjectFactory, minSize int, maxSize int) (err error) {
	value := p.pools.Get(name)
	if value != nil {
		return
	}
	ps, err := newPoolSet(minSize, maxSize, factory)
	if err != nil {
		return
	}
	p.pools.Set(name, ps)

	return

}
func (p *ObjectPool) UnRegister(name string) {
	obj := p.pools.Get(name)
	if obj == nil {
		return
	}
	p.pools.Delete(name)
	obj.(*poolSet).Close()
}

func (p *ObjectPool) Exists(name string) bool {
	return p.pools.Get(name) != nil
}

//Get 从对象组中申请一个对象
func (p *ObjectPool) Get(name string) (obj Object, err error) {
	v := p.pools.Get(name)
	if v == nil {
		err = errors.New(fmt.Sprint("not find pool: ", name))
		return
	}
	obj, err = v.(*poolSet).Get()
	if err != nil {
		err = fmt.Errorf("not find object from : %s,%s", name, err)
	}
	return

}

//Recycle 回收一个对象
func (p *ObjectPool) Recycle(name string, obj Object) {
	v := p.pools.Get(name)
	if v == nil {
		return
	}
	v.(*poolSet).back(obj)
}

//Unusable 标记为不可用，并通知连接池创建新的连接
func (p *ObjectPool) Unusable(name string, obj Object) {
	v := p.pools.Get(name)
	if v == nil {
		return
	}
	v.(*poolSet).createNew()
}

//Close 关闭一个对象
func (p *ObjectPool) close(name string) bool {
	v := p.pools.Get(name)
	if v == nil {
		return false
	}
	ps := v.(*poolSet)
	ps.Close()
	p.pools.Delete(name)
	return true

}
