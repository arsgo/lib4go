package concurrent

import (
	"reflect"
	"sync"
)

type CallBack func(...interface{}) (interface{}, error)

//ConcurrentMap 线程安全MAP
type ConcurrentMap struct {
	data map[string]interface{}
	lock *sync.RWMutex
}

//NewConcurrentMap 创建线程安全MAP
func NewConcurrentMap() (m *ConcurrentMap) {
	m = &ConcurrentMap{}
	m.lock = &sync.RWMutex{}
	m.data = make(map[string]interface{})
	return
}

//Exists 检查指定的KEY是否存在
func (c *ConcurrentMap) Exists(key string) (b bool) {
	c.lock.RLock()
	_, b = c.data[key]
	c.lock.RUnlock()
	return
}

//GetLength 获取总个数
func (c *ConcurrentMap) GetLength() (length int) {
	c.lock.RLock()
	length = len(c.data)
	c.lock.RUnlock()
	return
}

//Get 获取指定数据
func (c *ConcurrentMap) Get(key string) (data interface{}, b bool) {
	c.lock.RLock()
	data, b = c.data[key]
	c.lock.RUnlock()
	return
}

//Delete 删除指定KEY的数据
func (c *ConcurrentMap) Delete(key string) {
	c.lock.Lock()
	delete(c.data, key)
	c.lock.Unlock()
}

//Set 添加或修改指定KEY对应的值
func (c *ConcurrentMap) Set(key string, value interface{}) (b bool) {
	c.lock.Lock()
	b = !reflect.DeepEqual(c.data[key], value)
	c.data[key] = value
	c.lock.Unlock()
	return
}

//GetAll 获取所有数据
func (c *ConcurrentMap) GetAll() (nm map[string]interface{}) {
	nm = make(map[string]interface{})
	c.lock.RLock()
	for k, v := range c.data {
		nm[k] = v
	}
	c.lock.RUnlock()
	return
}

//GetAllAndClear 获取所有元素，并删除当前缓存列表
func (c *ConcurrentMap) GetAllAndClear() (nm map[string]interface{}) {
	nm = make(map[string]interface{})
	c.lock.RLock()
	for k, v := range c.data {
		nm[k] = v
	}
	c.lock.RUnlock()
	c.lock.Lock()
	c.data = make(map[string]interface{})
	c.lock.Unlock()
	return
}

//GetAndDel 获取当前值，并从缓存中删除
func (c *ConcurrentMap) GetAndDel(key string) (data interface{}, b bool) {
	c.lock.Lock()
	data, b = c.data[key]
	delete(c.data, key)
	c.lock.Unlock()
	return
}

//Clear 清除缓存中所有数据
func (c *ConcurrentMap) Clear() {
	c.lock.Lock()
	c.data = make(map[string]interface{})
	c.lock.Unlock()
	return
}

//GetOrAdd 添加或获取指定KEY
func (c *ConcurrentMap) GetOrAdd(key string, f CallBack, p ...interface{}) (b bool, data interface{}, err error) {
	c.lock.Lock()
	data, b = c.data[key]
	if b {
		b = false //未成功添加，数据已存在
		c.lock.Unlock()
		return
	}
	data, err = f(p...)
	if err != nil {
		c.lock.Unlock()
		return
	}
	b = true //数据成功添加
	c.data[key] = data
	c.lock.Unlock()
	return
}
