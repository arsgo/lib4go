package concurrent

import "fmt"

type keyValue struct {
	key   string
	value interface{}
}
type keyChan struct {
	key   string
	value chan interface{}
}
type keyChanMap struct {
	value chan map[string]interface{}
}

//ConcurrentMap 线程安全MAP
type ConcurrentMap struct {
	data     map[string]interface{}
	chSet    chan keyValue
	chDelete chan string
	chClose  chan int
	chGet    chan keyChan
	chAll    chan keyChanMap
}

//NewConcurrentMap 创建线程安全MAP
func NewConcurrentMap() (m ConcurrentMap) {
	m = ConcurrentMap{}
	m.data = make(map[string]interface{})
	m.chSet = make(chan keyValue, 1)
	m.chDelete = make(chan string, 1)
	m.chGet = make(chan keyChan, 1)
	m.chAll = make(chan keyChanMap, 1)
	go m.do()
	return
}

//Set 添加或修改指定KEY对应的值
func (c ConcurrentMap) Set(key string, value interface{}) {
	c.chSet <- keyValue{key: key, value: value}
}

//Delete 删除指定KEY的数据
func (c ConcurrentMap) Delete(key string) {
	c.chDelete <- key
}

//Get 获取指定KEY对应的数据
func (c ConcurrentMap) Get(key string) interface{} {
	ch := make(chan interface{}, 1)
	c.chGet <- keyChan{key: key, value: ch}
	value := <-ch
	return value
}

//GetAll 获取所有所有元素的拷贝
func (c ConcurrentMap) GetAll() map[string]interface{} {
	ch := make(chan map[string]interface{}, 1)
	c.chAll <- keyChanMap{value: ch}
	data := <-ch
	return data
}

//Close 关闭当前线程安全MAP
func (c ConcurrentMap) Close() {
	c.chClose <- 0
}

//do 单线程处理外部所有操作
func (c ConcurrentMap) do() {
LOOP:
	for {
		select {
		case data := <-c.chSet:
			{
				c.data[data.key] = data.value
			}
		case k := <-c.chDelete:
			fmt.Println("delete key:", k)
			delete(c.data, k)
		case g := <-c.chGet:
			{
				g.value <- c.data[g.key]

			}
		case ch := <-c.chAll:
			{
				data := make(map[string]interface{})
				for k, v := range c.data {
					data[k] = v
				}
				ch.value <- data
			}
		case <-c.chClose:
			break LOOP

		}
	}
}
