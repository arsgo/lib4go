package pool

import (
	"errors"
	"fmt"
	"sync"
)

//ObjectPool 对象缓存池
type ObjectPool struct {
	pools map[string]*poolSet
	mutex sync.Mutex
}

//New 创建一个新的对象k
func New() *ObjectPool {
	pools := &ObjectPool{}
	pools.pools = make(map[string]*poolSet)
	return pools
}

//Register 注册指定的对象组
func (p *ObjectPool) Register(name string, factory ObjectFactory, size int) int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[name]; ok {
		return v.list.Len()
	}

	if v, ok := p.pools[name]; ok {
		return v.list.Len()
	}
	p.pools[name] = newPoolSet(size, factory)
	return p.pools[name].list.Len()
}

//UnRegister 取消注册指定的对象组
func (p *ObjectPool) UnRegister(name string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[name]; ok {
		v.close()
		delete(p.pools, name)
	}

}
func (p *ObjectPool) Exists(name string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.pools[name]
	return ok
}

//Get 从对象组中申请一个对象
func (p *ObjectPool) Get(name string) (Object, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[name]; ok {
		return v.get()
	} else {
		return nil, errors.New(fmt.Sprintf("not find: %s", name))
	}
}

//Recycle 回收一个对象
func (p *ObjectPool) Recycle(name string, obj Object) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[name]; ok {
		v.back(obj)
	}
}

//Close 关闭一个对象组
func (p *ObjectPool) Close(name string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if ps, ok := p.pools[name]; ok && ps.usingCount == 0 {
		ps.close()
		delete(p.pools, name)
		return true
	}
	return false
}
