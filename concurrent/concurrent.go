package concurrent

import "strings"

const (
	GET = iota
	ADD
	SET
	DEL
	GETORADD
	ALL
	CLOSE
	LEN
)

type addResult struct {
	add   bool
	value interface{}
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
func (c *ConcurrentMap) Add(key string, f CallBack, p ...interface{}) (bool, interface{}) {
	if c.isClose || strings.EqualFold(key, "") {
		return false, nil
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, value: newswap(f, p...), method: ADD, result: ch}
	v := <-ch
	r := v.(*addResult)
	return r.add, r.value
}

//Set 添加或修改指定KEY对应的值
func (c *ConcurrentMap) Set(key string, value interface{}) {
	if c.isClose || strings.EqualFold(key, "") {
		return
	}
	c.request <- requestKeyValue{key: key, value: value, method: SET, result: make(chan interface{}, 1)}
}

//Delete 删除指定KEY的数据
func (c *ConcurrentMap) Delete(key string) {
	if c.isClose {
		return
	}
	c.request <- requestKeyValue{key: key, method: DEL}
}

//Get 获取指定KEY对应的数据
func (c *ConcurrentMap) Get(key string) interface{} {
	if c.isClose {
		return nil
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, method: GET, result: ch}
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
								data.result <- &addResult{add: false}
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
						c.data[data.key] = data.value
					}
				case CLOSE:
					c.isClose = true
				}
			}

		}
	}
}
