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
	canUse  int32
	added   int32
	notity  chan int
	isClose bool
}

//New 创建对象池
func newPoolSet(minSize int, maxSize int, fac ObjectFactory) (pool *poolSet, err error) {
	if maxSize == 0 {
		maxSize = 10
	}
	if minSize == 0 {
		minSize = 1
	}
	pool = &poolSet{minSize: minSize, maxSize: maxSize, factory: fac, queue: make(chan Object, maxSize),
		notity: make(chan int, maxSize)}
	go pool.init()
	return
}
func (p *poolSet) Get() (obj Object, err error) {
	if p.isClose {
		err = errors.New("cant get object from pool(pool is closed)")
		return
	}
	return p.getSingle(true)
}

func (p *poolSet) getSingle(create bool) (obj Object, err error) {
	defer func() {
		if atomic.LoadInt32(&p.canUse) == 0 && atomic.LoadInt32(&p.current) > 0 {
			p.createNew()
		}
	}()
	ticker := time.NewTicker(time.Millisecond * 50)
	select {
	case ps := <-p.queue:
		obj = ps
		atomic.AddInt32(&p.canUse, -1)
		return
	case <-ticker.C:
		break
	}
	err = fmt.Errorf("cant get object from pool:%d/%d/%d", atomic.LoadInt32(&p.current), p.minSize, p.maxSize)
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
func (p *poolSet) reCreate() {
	atomic.AddInt32(&p.current, -1)
	p.createNew()
}

func (p *poolSet) Close() {
	p.isClose = true
	for {
		obj, err := p.getSingle(false)
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
	if atomic.LoadInt32(&p.added) < int32(p.maxSize) {
		v := atomic.AddInt32(&p.added, 1)
		if v < int32(p.maxSize) {
			p.notity <- 1
		}
	}
}

//init 异步创建对象，factory.create要求返回正确可使用的对象，当对象不能创建成功时
// 该函数将持续堵塞，直到创建成功或收到关闭指定
func (p *poolSet) init() {
	for i := 0; i < p.minSize; i++ {
		p.createNew()
	}
	pk := time.NewTicker(time.Millisecond * 5)
	for {
		select {
		case _, ok := <-p.notity:
			if !ok || p.isClose {
				return
			}
			obj, err := p.factory.Create()
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * 5)
				p.createNew()
			} else {
				atomic.AddInt32(&p.current, 1)
				p.back(obj)
			}
		case <-pk.C:
			if p.isClose {
				return
			}
		}
	}

}
