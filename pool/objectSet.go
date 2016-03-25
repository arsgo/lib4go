package pool

import (
	"container/list"
	"sync"
)

type Object interface {
	Close()
}

//ObjectFactory
type ObjectFactory interface {
	Create() Object
}

//PoolSet
type poolSet struct {
	mutex   sync.Mutex
	Size    int
	list    *list.List
	factory ObjectFactory
}

//New 创建对象池
func newPoolSet(size int, fac ObjectFactory) *poolSet {
	pool := &poolSet{Size: size, factory: fac, list: list.New()}
	pool.init()
	return pool
}

func (p *poolSet) get() (Object, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	ele := p.list.Front()
	if ele == nil {
		object := p.factory.Create()
		return object, nil
	}
	p.list.Remove(ele)
	return ele.Value.(Object), nil
}

func (p *poolSet) add(obj Object) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.list.PushBack(obj)
}

func (p *poolSet) close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for p.list.Len()>0{
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
		p.list.PushBack(p.factory.Create())
	}
	return nil
}
