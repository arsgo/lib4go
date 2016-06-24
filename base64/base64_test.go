package base64

import "testing"

func TestAES(t *testing.T) {
	v := "abc34455=09"
	r := Encode(v)
	rv, err := Decode(r)
	if err != nil {
		t.Error(err)
	}
	if string(rv) != v {
		t.Error("解密失败：", string(rv), rv)
	}
}
