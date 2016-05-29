package concurrent

import (
	"strings"
	"testing"
)

var current ConcurrentMap

func init() {
	current = NewConcurrentMap()
}

func BenchmarkConcurrentMap(b *testing.B) {
	addChan := make(chan string, 10000)
	delChan := make(chan string, 10000)
	for i := 0; i < 10000; i++ {
		addChan <- "key_" + string(i)
	}
	go func() {
		for i := 0; i < 10000; i++ {
			select {
			case key := <-addChan:
				{
					current.Set(key, key)
					delChan <- key
				}

			}
		}

	}()
	for i := 0; i < 10000; i++ {
		select {
		case key := <-delChan:
			{
				v := current.Get(key)
				if strings.EqualFold(v.(string), "") {
					b.Error("get or set  data error")
				}
				current.Delete(key)
				v = current.Get(key)
				if v != nil {
					b.Error("del data error")
				}
			}
		}
	}
	data := current.GetAll()
	if len(data) != 0 {
		b.Error("batch delete error")
	}
	current.Close()

}

func TestConcurrentMap(t *testing.T) {
	value := current.Get("key")
	if value != nil {
		t.Error("get nil error")
	}
	all := current.GetAll()
	if len(all) != 0 {
		t.Error("get all error")
	}
	current.Set("key", "value")
	value = current.Get("key")
	if !strings.EqualFold(value.(string), "value") {
		t.Error("set value error")
	}
	all = current.GetAll()
	if len(all) != 1 {
		t.Error("get all error")
	}
	current.Delete("key")
	value = current.Get("key")
	if value != nil {
		t.Error("delete value error")
	}
	all = current.GetAll()
	if len(all) != 0 {
		t.Error("get all error")
	}

}
