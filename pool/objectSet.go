package pool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Object interface {
	Close()
}

//ObjectFactory
type ObjectFactory interface {
	Create() (Object, error)
	Close()
}

//PoolSet
type poolSet struct {
	mutex   sync.Mutex
	minSize int
	maxSize int
	queue   chan Object
	factory ObjectFactory
	current int32
	notity  chan int
	canUse  int32
	isClose bool
}

//New 创建对象池
func newPoolSet(minSize int, maxSize int, fac ObjectFactory) (pool *poolSet, err error) {
	pool = &poolSet{minSize: minSize, maxSize: maxSize, factory: fac, queue: make(chan Object, maxSize)}
	go pool.init()
	return
}

func (p *poolSet) get() (obj Object, err error) {
	if p.isClose || atomic.LoadInt32(&p.canUse) == 0 {
		err = errors.New("cant get object from pool")
		return
	}
	ticker := time.NewTicker(time.Millisecond * 50)
	select {
	case ps := <-p.queue:
		obj = ps
		atomic.AddInt32(&p.canUse, -1)
	case <-ticker.C:
		err = errors.New("cant get object from pool")
	}
	return
}

func (p *poolSet) back(obj Object) {
	if p.isClose {
		obj.Close()
		return
	}
	p.queue <- obj
	atomic.AddInt32(&p.canUse, 1)
}

func (p *poolSet) Close() {
	p.isClose = true
	for {
		obj, err := p.get()
		if err != nil {
			p.factory.Close()
			break
		} else {
			obj.Close()
		}
	}
}

//createNew 创建新的连接
func (p *poolSet) createNew() {
	if atomic.LoadInt32(&p.current) < int32(p.maxSize) {
		p.notity <- 1
	}
}

//init 异步创建对象，factory.create要求返回正确可使用的对象，当对象不能创建成功时
// 该函数将持续堵塞，直到创建成功或收到关闭指定
func (p *poolSet) init() {
	for i := 0; i < p.minSize; i++ {
		p.createNew()
	}
	pk := time.NewTicker(time.Millisecond * 50)
	for {
		select {
		case _, ok := <-p.notity:
			if !ok || p.isClose {
				return
			}
			obj, err := p.factory.Create()
			if err != nil {
				time.Sleep(time.Second * 5)
				p.createNew()
				fmt.Println("create failed:", err)
			} else {
				p.back(obj)
			}
		case <-pk.C:
			if p.isClose {
				return
			}
		}
	}

}
