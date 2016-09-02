package concurrent
/*
import (
	"reflect"
	"strings"
)


const (
	GET = iota
	ADD
	SET
	DEL
	EXISTS
	GETORADD
	GetANDDEL
	ALL
	ALLANDCLEAR
	CLEAR
	CLOSE
	LEN
)

type addResult struct {
	add   bool
	value interface{}
	err   error
}

type requestKeyValue struct {
	method int
	key    string
	value  interface{}
	result chan interface{}
}

type CallBack func(...interface{}) (interface{}, error)
type swap struct {
	params []interface{}
	call   CallBack
}

func newswap(call CallBack, p ...interface{}) *swap {
	return &swap{params: p, call: call}
}
func (s *swap) doCall() (interface{}, error) {
	return s.call(s.params...)
}

//ConcurrentMap 线程安全MAP
type ConcurrentMap struct {
	data    map[string]interface{}
	request chan requestKeyValue
	isClose bool
}

//NewConcurrentMap 创建线程安全MAP
func NewConcurrentMap() (m *ConcurrentMap) {
	m = &ConcurrentMap{isClose: false}
	m.data = make(map[string]interface{})
	m.request = make(chan requestKeyValue, 1000)
	go m.do()
	return
}

//Add 添加值
func (c *ConcurrentMap) Add(key string, f CallBack, p ...interface{}) (bool, interface{}, error) {
	if c.isClose || strings.EqualFold(key, "") {
		return false, false, nil
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, value: newswap(f, p...), method: ADD, result: ch}
	v := <-ch
	r := v.(*addResult)
	return r.add, r.value, r.err
}

//GetOrAdd 添加或获取指定KEY
func (c *ConcurrentMap) GetOrAdd(key string, f CallBack, p ...interface{}) (interface{}, error) {
	_, v, e := c.Add(key, f, p...)
	return v, e
}

//Set 添加或修改指定KEY对应的值
func (c *ConcurrentMap) Set(key string, value interface{}) bool {
	if c.isClose || strings.EqualFold(key, "") {
		return false
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, value: value, method: SET, result: ch}
	v := <-ch
	return v.(bool)
}

//Clear 清除所有数据
func (c *ConcurrentMap) Clear() {
	if c.isClose {
		return
	}
	c.request <- requestKeyValue{method: CLEAR}
}

//Delete 删除指定KEY的数据
func (c *ConcurrentMap) Delete(key string) {
	if c.isClose {
		return
	}
	c.request <- requestKeyValue{key: key, method: DEL}
}

//Exists 指定KEY是否存在
func (c *ConcurrentMap) Exists(key string) bool {
	if c.isClose {
		return false
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{result: ch, key: key, method: EXISTS}
	value := <-ch
	return value.(bool)
}

//Get 获取指定KEY对应的数据
func (c *ConcurrentMap) Get(key string) interface{} {
	if c.isClose {
		return nil
	}
	//start := time.Now()
	//	defer func() {
	//	tk := time.Now().Sub(start)
	//	if tk.Nanoseconds()/1000/1000 > 1 {
	//		fmt.Printf("+++++++end:%v\n", tk)
	//	}
	//	}()
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, method: GET, result: ch}
	value := <-ch
	return value
}

//GetAndDel 获取指定KEY对应的数据
func (c *ConcurrentMap) GetAndDel(key string) interface{} {
	if c.isClose {
		return nil
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, method: GetANDDEL, result: ch}
	value := <-ch
	return value
}

//GetLength 获取数据个数
func (c *ConcurrentMap) GetLength() int {
	if c.isClose {
		return 0
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{method: LEN, result: ch}
	value := <-ch
	return value.(int)
}

//GetAll 获取所有所有元素的拷贝
func (c *ConcurrentMap) GetAll() map[string]interface{} {
	if c.isClose {
		return make(map[string]interface{})
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{method: ALL, result: ch}
	data := <-ch
	if data != nil {
		return data.(map[string]interface{})
	}
	return make(map[string]interface{})
}

//GetAllAndClear 获取所有所有元素的拷贝
func (c *ConcurrentMap) GetAllAndClear() map[string]interface{} {
	if c.isClose {
		return make(map[string]interface{})
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{method: ALLANDCLEAR, result: ch}
	data := <-ch
	if data != nil {
		return data.(map[string]interface{})
	}
	return make(map[string]interface{})
}

//Close 关闭当前线程安全MAP
func (c *ConcurrentMap) Close() {
	if c.isClose {
		return
	}
	c.request <- requestKeyValue{method: CLOSE}
}

//do 单线程处理外部所有操作
func (c *ConcurrentMap) do() {
	for {
		select {
		case data := <-c.request:
			{
				switch data.method {
				case ADD:
					{
						if _, ok := c.data[data.key]; !ok {
							v, er := data.value.(*swap).doCall()
							if er != nil {
								data.result <- &addResult{add: false, err: er}
							} else {
								c.data[data.key] = v
								data.result <- &addResult{add: true, value: v}
							}
						} else {
							data.result <- &addResult{add: false, value: c.data[data.key]}
						}
					}
				case GET:
					{
						if d, ok := c.data[data.key]; ok {
							data.result <- d
						} else {
							data.result <- nil
						}
					}
				case EXISTS:
					{
						_, ok := c.data[data.key]
						data.result <- ok
					}
				case GetANDDEL:
					{
						if d, ok := c.data[data.key]; ok {
							data.result <- d
							delete(c.data, data.key)
						} else {
							data.result <- nil
						}
					}
				case ALLANDCLEAR:
					{
						values := make(map[string]interface{})
						for k, v := range c.data {
							values[k] = v
						}
						data.result <- values
						c.data = make(map[string]interface{})
					}
				case ALL:
					{
						values := make(map[string]interface{})
						for k, v := range c.data {
							values[k] = v
						}
						data.result <- values
					}
				case DEL:
					{
						delete(c.data, data.key)
					}
				case LEN:
					{
						data.result <- len(c.data)
					}
				case SET:
					{
						v := !reflect.DeepEqual(c.data[data.key], data.value)
						c.data[data.key] = data.value
						data.result <- v
					}
				case CLEAR:
					c.data = make(map[string]interface{})
				case CLOSE:
					c.isClose = true
				}
			}
		}
	}
}
*/