package pool

import (
    "sync"
    "errors"
    "fmt"
)

//ObjectPool 对象缓存池
type ObjectPool struct{
    pools map[string]*poolSet
    mutex sync.Mutex
}

//New 创建一个新的对象
func New()(*ObjectPool){
    pools:=&ObjectPool{}
    pools.pools=make(map[string]*poolSet)
    return  pools
}

//Register 注册指定的对象组
func (p *ObjectPool)Register(groupName string, factory ObjectFactory,size int)(int){
    if v,ok:=p.pools[groupName];ok{
        return v.list.Len()
    }
    p.mutex.Lock()
    defer p.mutex.Unlock()
     if v,ok:=p.pools[groupName];ok{
        return v.list.Len()
    }
    p.pools[groupName]=newPoolSet(size,factory)
    return p.pools[groupName].list.Len()
}

//Get 从对象组中申请一个对象
func (p *ObjectPool) Get(groupName string) (Object, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v,ok:=p.pools[groupName];ok{
       return v.get()
    }else{
        return nil,errors.New(fmt.Sprintf("%s not exists",groupName))
    }
}

//Recycle 回收一个对象
func (p *ObjectPool) Recycle(groupName string,obj Object) {
	p.mutex.Lock()
    defer p.mutex.Unlock()
    if v,ok:=p.pools[groupName];ok{
       v.add(obj)
    }
	
}

//Close 关闭一个对象组
func (p *ObjectPool) Close(groupName string){
    p.mutex.Lock()
    defer p.mutex.Unlock()
    if ps,ok:=p.pools[groupName];ok{
       ps.close()
       delete(p.pools,groupName)
    }  
}
