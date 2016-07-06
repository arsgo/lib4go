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

//ObjectPoolSnap 引擎池快照信息
type ObjectPoolSnap struct {
	Snaps []ObjectSnap `json:"objects"`
}

//ObjectSnap 引擎快照信息
type ObjectSnap struct {
	Name      string `json:"name"`
	Status    bool   `json:"status"`
	MinSize   int    `json:"minSize"`
	MaxSize   int    `json:"maxSize"`
	Available int    `json:"available"`
}

//New 创建一个新的对象k
func New() *ObjectPool {
	pools := &ObjectPool{}
	pools.pools = concurrent.NewConcurrentMap()
	return pools
}

//ResetAllPoolSize 重置所有链接池大小
func (p *ObjectPool) ResetAllPoolSize(minSize int, maxSize int) {
	if minSize == 0 || maxSize == 0 {
		err := fmt.Errorf("minSize:%d 或 maxSize:%d 不能为0", minSize, maxSize)
		fmt.Println(err)
		return
	}
	all := p.pools.GetAll()
	for _, value := range all {
		value.(*poolSet).SetSize(minSize, maxSize)
	}
}

//ResetPoolSize 重置链接池大小
func (p *ObjectPool) ResetPoolSize(name string, minSize int, maxSize int) {
	if minSize == 0 || maxSize == 0 {
		err := fmt.Errorf("%s minSize:%d 或 maxSize:%d 不能为0", name, minSize, maxSize)
		fmt.Println(err)
		return
	}
	value := p.pools.Get(name)
	if value != nil {
		value.(*poolSet).SetSize(minSize, maxSize)
	}
}

//Register 注册指定的对象
func (p *ObjectPool) Register(name string, factory ObjectFactory, minSize int, maxSize int) (err error) {
	if minSize == 0 || maxSize == 0 {
		err = fmt.Errorf("%s minSize:%d 或 maxSize:%d 不能为0", name, minSize, maxSize)
		fmt.Println(err)
		return
	}
	value := p.pools.Get(name)
	if value != nil {
		value.(*poolSet).SetSize(minSize, maxSize)
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
	v.(*poolSet).reCreate()
	obj.Close()
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

//GetSnap 获取当前连接池的快照信息
func (p *ObjectPool) GetSnap() (snaps ObjectPoolSnap) {
	snaps = ObjectPoolSnap{}
	pools := p.pools.GetAll()
	snaps.Snaps = make([]ObjectSnap, 0)
	for i, v := range pools {
		snap := ObjectSnap{}
		snap.Name = i
		set := v.(*poolSet)
		snap.Status = set.added > 0
		snap.MinSize = int(set.minSize)
		snap.MaxSize = int(set.maxSize)
		snap.Available = int(set.canUse)
		snaps.Snaps = append(snaps.Snaps, snap)
	}
	return snaps
}

//Close 关闭所有连接池
func (p *ObjectPool) Close() {
	all := p.pools.GetAll()
	for name := range all {
		p.close(name)
	}
}
