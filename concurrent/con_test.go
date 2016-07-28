package concurrent

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

var current ConcurrentMap

func init() {
	current = NewConcurrentMap()
}

func TestAdd(b *testing.B) {
	nmap := NewConcurrentMap()
	total := 0
	v := &sync.WaitGroup{}
	v.Add(3)
	f := func(p ...interface{}) (interface{}, error) {
		fmt.Println("A:", p[0])
		total++
		return p[0], errors.New("error")
	}

	go func() {
		for i := 0; i < 10000; i++ {
			nmap.Add(fmt.Sprintf("%d", i), f, i)
		}
		v.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			nmap.Add(fmt.Sprintf("%d", i), f, i)
		}
		v.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			nmap.Add(fmt.Sprintf("%d", i), f, i)
		}
		v.Done()
	}()
	v.Wait()
	if total != 10000 {
		b.Error("create count error")
	}
	if len(nmap.GetAll()) != 10000 {
		b.Error("added count error")
	}

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
