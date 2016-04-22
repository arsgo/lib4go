package pool

import (
	"container/list"
	"errors"
	"fmt"
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
func newPoolSet(size int, fac ObjectFactory) *poolSet {
	pool := &poolSet{Size: size, factory: fac, list: list.New()}
	pool.init()
	return pool
}

func (p *poolSet) get() (Object, error) {
	fmt.Println("get object from pool")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	ele := p.list.Front()
	if ele == nil {
		return nil, errors.New("cant get object")
	}
	p.list.Remove(ele)
	obj := ele.Value.(Object)
	if obj != nil {
		atomic.AddInt32(&p.usingCount, 1)
		return obj, nil
	}
	return nil, errors.New("no object can create")
}

func (p *poolSet) back(obj Object) {
	fmt.Println("back object from pool")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	atomic.AddInt32(&p.usingCount, -1)
	p.list.PushBack(obj)
}

func (p *poolSet) close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for p.list.Len() > 0 {
		ele := p.list.Front()
		if ele == nil {
			break
		}
		ele.Value.(Object).Close()
		p.list.Remove(ele)
	}
}

func (p *poolSet) init() error {
	for i := 0; i < p.Size; i++ {
		obj, err := p.factory.Create()
		if err != nil {
			panic(err)
		}
		p.list.PushBack(obj)
	}
	return nil
}
