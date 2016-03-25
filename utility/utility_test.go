package utility

import (
	"strings"
	"testing"
)

func Test_utility(t *testing.T) {

	if getLocalIP("127.0.0") != "127.0.0.1" {
		t.Error("get local ip error not equal 127.0.0.1")
	}
	if getLocalIP("192.168") != "192.168.0.245" {
		t.Error("get local ip error not equal 192.168.0.245")
	}
	if getLocalIP("1999") != "127.0.0.1" {
		t.Error("get local ip error")
	}

	data := NewDataMap()
	data.Set("ip", getLocalIP("127.0.0"))
	data.Set("domain", "/grs/delivery")
	result := data.Translate("@domain/bind.get/@ip")
	if !strings.EqualFold(result, "/grs/delivery/bind.get/127.0.0.1") {
		t.Errorf("translate error %s", result)
	}
    result = data.Translate("@domain/bind.get/{@ip}")
	if !strings.EqualFold(result, "/grs/delivery/bind.get/127.0.0.1") {
		t.Errorf("translate error %s", result)
	}
    

}
