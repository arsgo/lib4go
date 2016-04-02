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
func (p *ObjectPool) Register(groupName string, factory ObjectFactory, size int) int {
	if v, ok := p.pools[groupName]; ok {
		return v.list.Len()
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[groupName]; ok {
		return v.list.Len()
	}
	p.pools[groupName] = newPoolSet(size, factory)
	return p.pools[groupName].list.Len()
}

//UnRegister 取消注册指定的对象组
func (p *ObjectPool) UnRegister(groupName string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[groupName]; ok {
		v.close()
		delete(p.pools, groupName)
	}

}
func (p *ObjectPool) Exists(groupName string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.pools[groupName]
	return ok
}

//Get 从对象组中申请一个对象
func (p *ObjectPool) Get(groupName string) (Object, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[groupName]; ok {
		return v.get()
	} else {
		return nil, errors.New(fmt.Sprintf("not find: %s", groupName))
	}
}

//Recycle 回收一个对象
func (p *ObjectPool) Recycle(groupName string, obj Object) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.pools[groupName]; ok {
		v.back(obj)
	}

}

//Close 关闭一个对象组
func (p *ObjectPool) Close(groupName string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if ps, ok := p.pools[groupName]; ok && ps.usingCount == 0 {
		ps.close()
		delete(p.pools, groupName)
		return true
	}
	return false
}
