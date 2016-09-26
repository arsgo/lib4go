package script
/*
import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/pool"
)

//LuaPool  LUA对象池
type LuaPool struct {
	Binder  *LuaBinder
	vms     *concurrent.ConcurrentMap
	version int32
	minSize int
	maxSize int
	lk      sync.Mutex
	watcher *LuaScriptWatch
}

//NewLuaPool   构建LUA对象池
func NewLuaPool() *LuaPool {
	pool := &LuaPool{Binder: &LuaBinder{}, version: 0}
	pool.watcher = NewLuaScriptWatch(pool.Reload)
	pool.vms = concurrent.NewConcurrentMap()
	pool.vms.GetOrAdd(string(pool.version+1), pool.createVM, pool.version, pool.version+1)
	return pool
}

//SetPoolSize 设置连接池大小
func (p *LuaPool) SetPoolSize(minSize int, maxSize int) {
	p.minSize = minSize
	p.maxSize = maxSize
}

//Call 选取最新的脚本引擎执行当前脚本
func (p *LuaPool) Call(input InputArgs) (result []string, outparams map[string]string, err error) {
	var ok bool
	input.Content, ok = p.watcher.GetContent(input.Script)
	if !ok {
		err = fmt.Errorf("脚本不存在:%s", input.Script)
		return
	}
	return p.getVM().Call(input)
}

//Reload 重新加载所有引擎
func (p *LuaPool) Reload() {
	current := atomic.LoadInt32(&p.version)
	next := current + 1
	if b, _, er := p.vms.GetOrAdd(string(next), p.createVM, current, next); b && er == nil {
		lastVM, ok := p.vms.Get(string(current))
		if ok {
			p.vms.Delete(string(current))
			lastVM.(*luavm).Close()
		}
	}
}

//PreLoad 预加载脚本
func (p *LuaPool) PreLoad(script string, minSize int, maxSize int) error {
	p.Binder.RegisterScript(script, minSize, maxSize)
	return p.watcher.AppendFile(script)
}

//GetSnap 获取LUA引擎的快照信息
func (p *LuaPool) GetSnap() pool.ObjectPoolSnap {
	return p.getVM().GetSnap()
}

//Close 关闭引擎
func (p *LuaPool) Close() {
	p.getVM().Close()
}

func (p *LuaPool) getVM() *luavm {
	vm, _ := p.vms.Get(string(p.version))
	return vm.(*luavm)
}
func (p *LuaPool) createVM(args ...interface{}) (interface{}, error) {
	p.lk.Lock()
	defer p.lk.Unlock()
	lastVesion := args[0].(int32)
	currentVersion := args[1].(int32)
	if atomic.LoadInt32(&p.version) != lastVesion {
		return nil, errors.New("创建失败，版本错误")
	}
	vm := newLuaVM(p.Binder, p.minSize, p.maxSize)
	if atomic.CompareAndSwapInt32(&p.version, lastVesion, currentVersion) {
		return vm, nil
	}
	return nil, errors.New("创建失败，版本错误")
}
*/