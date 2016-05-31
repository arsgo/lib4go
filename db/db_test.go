package db

import (
	"strings"
	"testing"
)

func TestBaseQuerySelect(t *testing.T) {
	orcl, err := NewDB("oracle", "grs_delivery/123456@ORCL136")
	if err != nil {
		t.Error(err)
	}
	orcl.SetPoolSize(5, 10)
	data, err := orcl.Query("select 'colin' name, 1 id from dual")
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Error("查询返回数据条数有误")
	}
	name := data[0]["NAME"]
	for i, v := range data[0] {
		t.Log(i, v)
	}
	if name == nil || !strings.EqualFold(name.(string), "colin") {
		t.Error("返回结果有误:", len(data), name)
	}
}
func TestFromDb(t *testing.T) {
	orcl, err := NewDB("oracle", "grs_delivery/123456@ORCL136")
	if err != nil {
		t.Error(err)
	}
	orcl.SetPoolSize(5, 10)
	data, err := orcl.Query("select to_char(sysdate,'yyyymmddhh24miss') time from dual")
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Error("查询返回数据条数有误")
	}
	for i, v := range data[0] {
		t.Log(i, v)
	}
}
func TestProcedure(t *testing.T) {
	orcl, err := NewDB("oracle", "grs_delivery/123456@ORCL136")
	if err != nil {
		t.Error(err)
	}
	orcl.SetPoolSize(5, 10)
	data, err := orcl.Query("call gr_p_mobile_get(?)", "abc")
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Error("查询返回数据条数有误")
	}
	for i, v := range data[0] {
		t.Log(i, v)
	}
}
