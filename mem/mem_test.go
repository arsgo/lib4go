package mem

import (
	"strings"
	"testing"
)

func Test_memcache(t *testing.T) {
	key := "t1_key"
	value := "t1_value1231231"
	mem := New("192.168.101.161:11211")
	if e:=mem.Set(key, value, 3000);e!= nil {
		t.Error(e)
	}
	v := mem.Get(key)	
	if !strings.EqualFold(v, value) {
		t.Errorf("get value error [%s]",v)
	}
	mem.Delete(key)

	v = mem.Get(key)	
	if !strings.EqualFold(v, "") {
		t.Errorf("get value error [%s]",v)
	}

}
