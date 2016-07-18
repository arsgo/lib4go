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
)

type requestKeyValue struct {
	method int
	key    string
	value  interface{}
	result chan interface{}
}

//ConcurrentMap 线程安全MAP
type ConcurrentMap struct {
	data    map[string]interface{}
	request chan requestKeyValue
	isClose bool
}

//NewConcurrentMap 创建线程安全MAP
func NewConcurrentMap() (m ConcurrentMap) {
	m = ConcurrentMap{}
	m.data = make(map[string]interface{})
	m.request = make(chan requestKeyValue, 100)
	go m.do()
	return
}

//Add 添加值
func (c ConcurrentMap) Add(key string, value interface{}) bool {
	if c.isClose || strings.EqualFold(key, "") {
		return false
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, value: value, method: ADD, result: ch}
	v := <-ch
	return v.(bool)
}

//Set 添加或修改指定KEY对应的值
func (c ConcurrentMap) Set(key string, value interface{}) {
	if c.isClose || strings.EqualFold(key, "") {
		return
	}
	c.request <- requestKeyValue{key: key, value: value, method: SET, result: make(chan interface{}, 1)}
}

//Delete 删除指定KEY的数据
func (c ConcurrentMap) Delete(key string) {
	if c.isClose {
		return
	}
	c.request <- requestKeyValue{key: key, method: DEL}
}

//GetOrAdd 不存在时添加，存在时直接返回值
func (c ConcurrentMap) GetOrAdd(key string, value interface{}) interface{} {
	if c.isClose || strings.EqualFold(key, "") {
		return nil
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, value: value, method: GETORADD, result: ch}
	v := <-ch
	return v
}

//Get 获取指定KEY对应的数据
func (c ConcurrentMap) Get(key string) interface{} {
	if c.isClose {
		return nil
	}
	ch := make(chan interface{}, 1)
	c.request <- requestKeyValue{key: key, method: GET, result: ch}
	value := <-ch
	return value
}

//GetLength 获取数据个数
func (c ConcurrentMap) GetLength() int {
	return len(c.data)
}

//GetAll 获取所有所有元素的拷贝
func (c ConcurrentMap) GetAll() map[string]interface{} {
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
func (c ConcurrentMap) Close() {
	if c.isClose {
		return
	}
	c.request <- requestKeyValue{method: CLOSE}
}

//do 单线程处理外部所有操作
func (c ConcurrentMap) do() {
	for {
		select {
		case data := <-c.request:
			{
				switch data.method {
				case ADD:
					{
						if _, ok := c.data[data.key]; !ok {
							c.data[data.key] = data.value
							data.result <- true
						} else {
							data.result <- false
						}
					}
				case GETORADD:
					{
						if value, ok := c.data[data.key]; !ok {
							c.data[data.key] = data.value
							data.result <- data.value
						} else {
							data.result <- value
						}
					}
				case GET:
					{
						data.result <- c.data[data.key]
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
