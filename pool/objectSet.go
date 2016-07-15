package pool
/*
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
	minSize int32
	maxSize int32
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
	pool = &poolSet{minSize: int32(minSize), maxSize: int32(swap(maxSize, 999)), factory: fac}
	pool.queue = make(chan Object, pool.maxSize*10)
	pool.notity = make(chan int, pool.maxSize*10)
	go pool.init()
	return
}
func (p *poolSet) SetSize(min int, max int) {
	p.resetMinMaxSize(min, max)
}
func (p *poolSet) resetMinMaxSize(min int, max int) {
	atomic.SwapInt32(&p.maxSize, int32(max))
	minValue := atomic.SwapInt32(&p.minSize, int32(min))
	remain := int(atomic.LoadInt32(&p.minSize) - minValue)
	if remain > 0 {
		fmt.Printf("reset pool set min:%d,max:%d\r\n", min, max)
	}
	for i := 0; i < remain; i++ {
		p.createNew()
	}
}
func (p *poolSet) Get() (obj Object, err error) {
	if p.isClose {
		err = errors.New("cant get object from pool(pool is closed)")
		return
	}
	return p.getSingle(true)
}

func (p *poolSet) getSingle(create bool) (obj Object, err error) {
	if atomic.LoadInt32(&p.canUse) == 0 {
		p.createNew()
	}
	timeOut := time.NewTicker(time.Millisecond * 80)
	createNew := time.NewTicker(time.Millisecond * 41)
BRK:
	for {
		select {
		case ps := <-p.queue:
			obj = ps
			atomic.AddInt32(&p.canUse, -1)
			break BRK
		case <-createNew.C:
			p.createNew()
		case <-timeOut.C:
			err = fmt.Errorf("cant get object from pool:%d/%d/%d/%d", atomic.LoadInt32(&p.canUse), atomic.LoadInt32(&p.current), p.minSize, p.maxSize)
			break BRK
		}
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
	for i := 0; i < int(atomic.LoadInt32(&p.minSize)); i++ {
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
func swap(v int, def int) int {
	if v == 0 {
		return def
	}
	return v
}
*/