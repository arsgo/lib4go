package pool

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/arsgo/lib4go/concurrent"
)

//ObjectPool 对象缓存池, 缓存池中的对象只添加,不会修改或删除,部分代码对锁进行了优化
type ObjectPool struct {
	pools   *concurrent.ConcurrentMap
	using   int32
	isClose int32
}

//ObjectPoolSnap 引擎池快照信息
type ObjectPoolSnap struct {
	Snaps []ObjectSnap `json:"objects"`
}

//ObjectSnap 引擎快照信息
type ObjectSnap struct {
	Name    string `json:"name"`
	Status  bool   `json:"status"`
	MinSize int    `json:"minSize"`
	MaxSize int    `json:"maxSize"`
	Cache   int    `json:"cache"`
}

//New 创建一个新的对象k
func New() *ObjectPool {
	pools := &ObjectPool{isClose: 1, using: 0}
	pools.pools = concurrent.NewConcurrentMap()
	return pools
}

func (o *ObjectPool) createSet(p ...interface{}) (interface{}, error) {
	minSize, maxSize, factory := p[0].(int), p[1].(int), p[2].(ObjectFactory)
	return newPoolSet(minSize, maxSize, factory)
}

//Register 注册指定的对象
func (o *ObjectPool) Register(name string, factory ObjectFactory, minSize int, maxSize int) (err error) {
	if atomic.LoadInt32(&o.isClose) == 0 {
		err = errors.New(fmt.Sprint("pool is closed", name))
		return
	}
	o.pools.Add(name, o.createSet, minSize, maxSize, factory)
	return nil
}

//Get 从对象组中申请一个对象
func (o *ObjectPool) Get(name string) (obj Object, err error) {
	if atomic.LoadInt32(&o.isClose) == 0 {
		err = errors.New(fmt.Sprint("pool is closed", name))
		return
	}
	v := o.pools.Get(name)
	if v == nil {
		err = errors.New(fmt.Sprint("not find pool: ", name))
		return
	}
	obj, err = v.(*poolSet).Get()
	if err != nil {
		err = fmt.Errorf("not find object from : %s,%s", name, err)
	}
	atomic.AddInt32(&o.using, 1)
	return

}

//Recycle 回收一个对象
func (o *ObjectPool) Recycle(name string, obj Object) {
	atomic.AddInt32(&o.using, -1)
	v := o.pools.Get(name)
	if v == nil {
		return
	}
	v.(*poolSet).Back(obj)
}

//Unusable 标记为不可用，并通知连接池创建新的连接
func (o *ObjectPool) Unusable(name string, obj Object) {
	atomic.AddInt32(&o.using, -1)
	v := o.pools.Get(name)
	if v == nil {
		return
	}
	v.(*poolSet).startCacheMaker()
	obj.Close()
}

//UnRegister 关闭一个对象
func (o *ObjectPool) UnRegister(name string) bool {
	v := o.pools.Get(name)
	if v == nil {
		return false
	}
	ps := v.(*poolSet)
	ps.Close()
	o.pools.Delete(name)
	return true

}

//GetSnap 获取当前连接池的快照信息
func (o *ObjectPool) GetSnap() (snaps ObjectPoolSnap) {
	snaps = ObjectPoolSnap{}
	pools := o.pools.GetAll()
	snaps.Snaps = make([]ObjectSnap, 0)
	for i, v := range pools {
		snap := ObjectSnap{}
		snap.Name = i
		set := v.(*poolSet)
		snap.Status = set.hasCratedCount > 0
		snap.MinSize = int(set.minSize)
		snap.MaxSize = int(set.maxSize)
		snap.Cache = int(set.cacheCount)
		snaps.Snaps = append(snaps.Snaps, snap)
	}
	return snaps
}

//Close 关闭所有连接池
func (o *ObjectPool) Close() {
	if atomic.CompareAndSwapInt32(&o.isClose, 1, 0) {
		go o.startClose()
	}
}
func (o *ObjectPool) startClose() {
	timeout := time.NewTicker(time.Second * 31)
	timeChecker := time.NewTicker(time.Second * 2)
START:
	for {
		select {
		case <-timeout.C:
			break START
		case <-timeChecker.C:
			if atomic.LoadInt32(&o.using) > 0 {
				continue
			}
			break START
		}
	}
	all := o.pools.GetAll()
	for name := range all {
		//	fmt.Println("-----关闭引擎:", name)
		o.UnRegister(name)
	}
}

func (o *ObjectPool) Exists(name string) bool {
	return o.pools.Get(name) != nil
}
