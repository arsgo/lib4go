package aes

import (
	"strings"
	"testing"
)

func TestAES(t *testing.T) {
	key := "1234567890123456"
	msg := "abc123!?$*&()'-=@~"
	sec, err := Encrypt(msg, key)
	if err != nil {
		t.Error(err)
	}
	nmsg, err := Decrypt(sec, key)
	if err != nil {
		t.Error(err)
	}
	if !strings.EqualFold(nmsg, msg) {
		t.Errorf("解密失败:[%s][%s]", nmsg, msg)
	}
	t.Log(sec)
	t.Log(nmsg)

}
