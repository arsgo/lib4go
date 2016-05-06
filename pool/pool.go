package pool

import (
	"errors"
	"fmt"
	"sync"
)

//ObjectPool 对象缓存池, 缓存池中的对象只添加,不会修改或删除,部分代码对锁进行了优化
type ObjectPool struct {
	pools map[string]*poolSet
	//mutex sync.Mutex
	lock sync.RWMutex
}

//New 创建一个新的对象k
func New() *ObjectPool {
	pools := &ObjectPool{}
	pools.pools = make(map[string]*poolSet)
	return pools
}

//Register 注册指定的对象
func (p *ObjectPool) Register(name string, factory ObjectFactory, size int) (err error) {
	if _, ok := p.pools[name]; ok {
		return
	}
	p.lock.Lock()
	if _, ok := p.pools[name]; !ok {
		ps, err := newPoolSet(size, factory)
		if err == nil {
			p.pools[name] = ps
		}
	}
	p.lock.Unlock()
	return

}
func (p *ObjectPool) UnRegister(name string) {
	p.lock.Lock()
	if v, ok := p.pools[name]; ok {
		v.close()
		delete(p.pools, name)
	}
	p.lock.Unlock()

}

func (p *ObjectPool) Exists(name string) bool {
	p.lock.RLock()
	_, ok := p.pools[name]
	p.lock.RUnlock()
	return ok
}

//Get 从对象组中申请一个对象
func (p *ObjectPool) Get(name string) (obj Object, err error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if v, ok := p.pools[name]; ok {
		obj, err = v.get()
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("not find: %s from pool", name))
	}
	return

}

//Recycle 回收一个对象
func (p *ObjectPool) Recycle(name string, obj Object) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if v, ok := p.pools[name]; ok {
		v.back(obj)
	}
}

//Close 关闭一个对象
func (p *ObjectPool) Close(name string) bool {
	p.lock.Lock()
	if ps, ok := p.pools[name]; ok && ps.usingCount == 0 {
		ps.close()
		delete(p.pools, name)
		p.lock.Lock()
		return true
	}
	p.lock.Lock()
	return false
}
