/*
1. 对象池管理，初始化时当最大缓存数量为0时使用最小缓存数量，当最小缓存0时，不启动缓存创建线程
2. 当最小缓存数量大于0时，启动缓存创建线程，并创建数量达到应创建数量时，退出该线程，创建失败时定时重新创建
3. Get从缓存区获取对象，当获取失败后立即重新创建对象，获取成功则直接返回
4. 创建新对象时发生错误则启动缓存创建线程创建对象

*/

package pool

import (
	"fmt"
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

type poolSet struct {
	name               string
	minSize            int32       //最小缓存个数
	maxSize            int32       //最大缓存个数
	createTimes        int32       //创建次数
	isCreatingStatus   int32       //是否正在创建
	isCreatingCount    int32       //正在创建个数
	hasCratedCount     int32       //已经创建个数
	cacheCount         int32       //当前缓存可用数
	getCount           int32       //当前已取出个数
	createMessageQueue chan int    //创建消息
	isClose            int32       //是否已关闭
	closeMessageQueue  chan int    //关闭消息
	cacheQueue         chan Object //可用缓存
	factory            ObjectFactory
	lastUseTime        time.Time //最后使用时间
	timeout            time.Duration
}

//newPoolSet 创建对象管理器，当maxSize<=0时，设置默认值为100个，当minSize<=0时默认值为0
//当minSize==0时，不创建缓存对象，当minSize>0时启动新的goroutine创建当minSize个数的对象，当创建成功后退出该当goroutine
//创建失败后不退出该goroutine,并定时创建该Object
//当调用Get函数创建新的object失败后，会检查是否启动Watch goroutine创建对象，如果未创建则创建goroutine定时创建对象，
//此时会检查minSize是否为0，为0则会加入创建消息，于Watch创建
func newPoolSet(name string, minSize int, maxSize int, fac ObjectFactory) (pool *poolSet, err error) {
	pool = &poolSet{name: name, minSize: int32(minSize), maxSize: int32(swap(maxSize, swap(minSize, 500))), factory: fac, createTimes: 2}
	pool.timeout = time.Second * 10
	pool.createMessageQueue = make(chan int, pool.maxSize)
	pool.cacheQueue = make(chan Object, pool.maxSize)
	pool.closeMessageQueue = make(chan int, 1)
	pool.startInit()
	go pool.clear()
	return

}

//createNew 创建新的Object，创建失败时启动watch流程用于定时创建
//成功时累加已创建个数
func (p *poolSet) createNew() (obj Object, err error) {
	if obj, err = p.factory.Create(); err != nil {
		go p.startCacheMaker(0)
		return
	}
	atomic.AddInt32(&p.hasCratedCount, 1)
	go p.startCacheMaker(1)
	return
}

//createOne 创建新的Object，创建失败时启动watch流程用于定时创建
//成功时累加已创建个数
func (p *poolSet) createOne() (obj Object, err error) {
	if obj, err = p.factory.Create(); err != nil {
		return
	}
	atomic.AddInt32(&p.hasCratedCount, 1)
	return
}

func (p *poolSet) Close() {
	p.factory.Close()
	if atomic.CompareAndSwapInt32(&p.isClose, 0, 1) {
		p.closeMessageQueue <- 1
	START:
		for {
			select {
			case ps := <-p.cacheQueue:
				atomic.AddInt32(&p.cacheCount, -1)
				ps.Close()
			default:
				break START
			}
		}
	}
}

//Get 获取可用的Object,当没有可用时立即创建新的Object
func (p *poolSet) Get() (obj Object, err error) {
	if atomic.LoadInt32(&p.isClose) != 0 {
		err = fmt.Errorf("pool is closed")
		return
	}
	atomic.AddInt32(&p.getCount, 1)
	select {
	case obj = <-p.cacheQueue:
		atomic.AddInt32(&p.cacheCount, -1)
	default:
	}
	if obj == nil {
		obj, err = p.createNew()
	}
	p.lastUseTime = time.Now()
	return
}

//Back 回收Object当，未超过最大缓存大小时回收对象，否则将关闭object并丢弃
func (p *poolSet) Back(obj Object) {
	if p.isCreatingStatus == 1 {
		atomic.AddInt32(&p.cacheCount, 1)
		select {
		case p.cacheQueue <- obj:
		default:
			obj.Close()
		}
		return
	}
	if atomic.AddInt32(&p.cacheCount, 1) <= atomic.LoadInt32(&p.maxSize) {
		select {
		case p.cacheQueue <- obj:
		default:
			obj.Close()
		}
	} else {
		atomic.AddInt32(&p.cacheCount, -1)
		obj.Close()
	}
}

//startInit 启动初始化
func (p *poolSet) startInit() {
	if atomic.LoadInt32(&p.minSize) == 0 {
		return
	}
	for i := 0; i < int(atomic.LoadInt32(&p.minSize)-atomic.LoadInt32(&p.isCreatingCount)-
		atomic.LoadInt32(&p.hasCratedCount))+i; i++ {
		atomic.AddInt32(&p.isCreatingCount, 1)
		select {
		case p.createMessageQueue <- i:
		default:
		}
	}
	//修改创建状态
	if atomic.CompareAndSwapInt32(&p.isCreatingStatus, 0, 1) {
		go p.makeCache()
	}
}

//startCacheMaker,检查是否需要创建监控程序
//1. 已创建个数小于最大缓存个数
//2. 创建状态为未启动
func (p *poolSet) startCacheMaker(r int32) {
	if r < 0 {
		return
	}
	if !atomic.CompareAndSwapInt32(&p.isCreatingStatus, 0, 1) {
		return
	}
	c := p.createTimes*2 - r
	atomic.AddInt32(&p.createTimes, 1)
	count := min(c, p.maxSize-p.cacheCount)
	//添加创建消息
	if atomic.CompareAndSwapInt32(&p.isCreatingCount, 0, count) {
		for i := 0; i < int(count); i++ {
			select {
			case p.createMessageQueue <- -1:
			}

		}
	}
	go p.makeCache()

}

//makeCache 获取创建消息，并循环创建对象，当全部创建成功后自动退出
func (p *poolSet) makeCache() {
START:
	for {
		select {
		case <-p.closeMessageQueue:
			break START
		case i := <-p.createMessageQueue:
			obj, err := p.createOne() //创建新的Object
			if err != nil {
				fmt.Println(err)
				p.createMessageQueue <- i
				time.Sleep(p.timeout) //创建失败休息一定时间继续创建
				continue
			}
			p.push(obj)
			if atomic.AddInt32(&p.isCreatingCount, -1) == 0 { //当正在创建数小于0时退出循环
				break START
			}
		}
	}
	atomic.CompareAndSwapInt32(&p.isCreatingStatus, 1, 0) //切换状态
}

//push 将新创建的Object放入缓存
func (p *poolSet) push(obj Object) bool {
	if atomic.AddInt32(&p.cacheCount, 1) <= atomic.LoadInt32(&p.maxSize) {
		select {
		case p.cacheQueue <- obj:
			return true
		default:
			atomic.AddInt32(&p.cacheCount, -1)
			obj.Close()
		}

	} else {
		atomic.AddInt32(&p.cacheCount, -1)
		obj.Close()
	}
	return false
}

//clear 减少缓存数量
func (p *poolSet) clear() {
	tk := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-tk.C:
		START:
			//	fmt.Printf("%s已缓存对象:%d-%d\n", p.name, p.cacheCount, len(p.cacheQueue))
			if int(p.getCount) < swap(int(p.minSize), 1)*100 && p.cacheCount > p.minSize {
				select {
				case obj := <-p.cacheQueue:
					atomic.AddInt32(&p.cacheCount, -1)
					atomic.StoreInt32(&p.createTimes, 1)
					obj.Close()
					goto START
				default:
				}
			}
			atomic.StoreInt32(&p.getCount, 0)
		}
	}
}

func swap(v int, def int) int {
	if v == 0 {
		return def
	}
	return v
}
func min(a int32, b int32) int32 {
	if a > b {
		return b
	}
	return a
}
