package pool

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

type Object interface {
	Close()
	Check() bool
	Fatal()
}

//ObjectFactory
type ObjectFactory interface {
	Create() (Object, error)
}

//PoolSet
type poolSet struct {
	mutex      sync.Mutex
	Size       int
	list       *list.List
	factory    ObjectFactory
	usingCount int32
}

//New 创建对象池
func newPoolSet(size int, fac ObjectFactory) (pool *poolSet, err error) {
	pool = &poolSet{Size: size, factory: fac, list: list.New()}
	go pool.init()
	return
}

func (p *poolSet) get() (obj Object, err error) {
	p.mutex.Lock()
	ele := p.list.Front()
	if ele != nil {
		p.list.Remove(ele)
		p.mutex.Unlock()
		obj = ele.Value.(Object)
		if obj != nil {
			atomic.AddInt32(&p.usingCount, 1)
			return
		}
	} else {
		p.mutex.Unlock()
	}
	err = errors.New("cant get object from pool")
	return
}

func (p *poolSet) back(obj Object) {
	p.mutex.Lock()
	p.list.PushBack(obj)
	p.mutex.Unlock()
	atomic.AddInt32(&p.usingCount, -1)
}

func (p *poolSet) close() {
	p.mutex.Lock()
	var ele *list.Element
	for p.list.Len() > 0 {
		ele = p.list.Front()
	}
	p.mutex.Unlock()
	if ele != nil {
		ele.Value.(Object).Close()
		p.list.Remove(ele)
	}

}

func (p *poolSet) init() error {
	for i := 0; i < p.Size; i++ {
		obj, err := p.factory.Create()
		if err != nil {
			return err
		}
		p.mutex.Lock()
		p.list.PushBack(obj)
		p.mutex.Unlock()
	}
	return nil
}
