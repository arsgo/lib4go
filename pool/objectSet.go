package pool

import (
	"container/list"
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
	mutex      sync.Mutex
	Size       int
	list       *list.List
	factory    ObjectFactory
	usingCount int32
	makeQueue  chan int
	isClose    bool
}

//New 创建对象池
func newPoolSet(size int, fac ObjectFactory) (pool *poolSet, err error) {
	pool = &poolSet{Size: size, factory: fac, list: list.New()}
	pool.makeQueue = make(chan int, size)
	go pool.init()
	return
}

func (p *poolSet) get() (obj Object, err error) {
	if p.isClose {
		err = errors.New("cant get object from pool")
		return
	}
	fmt.Println("get->get lock")
	p.mutex.Lock()
	defer func() {
		p.mutex.Unlock()
		fmt.Println("get->back lock")
	}()
	ele := p.list.Front()
	if ele != nil {
		p.list.Remove(ele)
		obj = ele.Value.(Object)
		if obj != nil {
			atomic.AddInt32(&p.usingCount, 1)
			return
		}
	}
	err = errors.New("cant get object from pool")
	return
}

func (p *poolSet) back(obj Object) {
	if p.isClose {
		obj.Close()
		return
	}
	fmt.Println("back->get lock")
	p.mutex.Lock()
	defer func() {
		p.mutex.Unlock()
		fmt.Println("back->back lock")
	}()
	p.list.PushBack(obj)
	atomic.AddInt32(&p.usingCount, -1)
}

func (p *poolSet) close() {
		p.isClose = true
	fmt.Println("close current poolset")
	fmt.Println("close->get lock")

	p.mutex.Lock()
	defer func() {
		p.mutex.Unlock()
		fmt.Println("close->back lock")
	}()
	var ele *list.Element
	for p.list.Len() > 0 {
		ele = p.list.Front()
	}

	if ele != nil {
		fmt.Println("start close object")
		ele.Value.(Object).Close()
		fmt.Println("end close object")
		p.list.Remove(ele)
	}
	fmt.Println("start close factory")
	p.factory.Close()
	fmt.Println("end close factory")
	

}

//createNew 创建新的连接
func (p *poolSet) createNew() {
	p.makeQueue <- -1
}

//init 异步创建对象，factory.create要求返回正确可使用的对象，当对象不能创建成功时
// 该函数将持续堵塞，直到创建成功或收到关闭指定
func (p *poolSet) init() {
	for i := 0; i < p.Size; i++ {
		p.makeQueue <- i
	}
	pk := time.NewTicker(time.Second)
	defer pk.Stop()
	for {
		select {
		case v, ok := <-p.makeQueue:
			if !ok || p.isClose {
				return
			}
			obj, err := p.factory.Create()
			if err != nil {
				time.Sleep(time.Second * 10)
				p.makeQueue <- v
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
