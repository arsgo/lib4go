package des

import (
	"strings"
	"testing"
)

func TestDes(t *testing.T) {
	key := "12345678"
	value := "987654321"
	expect := "97c465b54478cd538bc8b61c4e75b7a6bfd7c6d1bc4ead67"
	actual, err := Encrypt(value, key)
	if err != nil || !strings.EqualFold(expect, actual) {
		t.Error("加密结果错误",err)
	}
	original, err := Decrypt(actual, key)
	if err != nil || !strings.EqualFold(original, expect) {
		t.Error("解密结果错误",err)
	}
}
