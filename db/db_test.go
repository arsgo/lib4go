package db

import "testing"
import "fmt"
import "encoding/json"
import "github.com/arsgo/lib4go/utility"

/*
func TestScalar(t *testing.T) {
	t.Log("sales")
	orcl, err := NewDBMapByConfig(`{
    "provider":"oracle",
    "connString":"CY_ESALES/123456@ORCL136"
}`)
	if err != nil {
		t.Error(err)
	}
	data, err := orcl.Query(`select t.product_name,t.product_id
	from
	cy_product_info t`, make(map[string]interface{}))
	if err != nil {
		t.Error(err)
	}
	t.Log(data.Result)
}
*/
/*
func TestBaseQuerySelect(t *testing.T) {
	orcl, err := NewDB("oracle", "grs_delivery/123456@ORCL136")
	if err != nil {
		t.Error(err)
	}
	orcl.SetPoolSize(5, 10)
	data, err := orcl.Query("select 'colin' name, 1 id from dual where 1=:1", 1)
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


func TestSales(t *testing.T) {
	t.Log("sales")
	orcl, err := NewDBMapByConfig(`{
    "provider":"oracle",
    "connString":"CY_ESALES/123456@ORCL136"
}`)
	if err != nil {
		t.Error(err)
	}
	data, err := orcl.Query(`select t.product_name
	from
	cy_product_info t`, make(map[string]interface{}))
	if err != nil {
		t.Error(err)
	}
	t.Log(data.Result)
	buffer, err := json.Marshal(&data)
	t.Log(string(buffer), err)
	if len(data.Result) != 1 {
		t.Error("查询返回数据条数有误")
	}
	for i, v := range data.Result {
		t.Log(i, "[", v, "]")
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
	tr, err := orcl.Begin()
	if err != nil {
		t.Error(err)
	}
	row, err = tr.Execute("update gr_order_notify t set t.notify_now=:1 where t.notify_id=:2", 9, 27277522)
	if err != nil {
		tr.Rollback()
		t.Error(err)
	}
	if row != 1 {
		t.Error("返回条数或结果不正确")
	}
	tr.Commit()

}
func TestSchema(t *testing.T) {
	input := map[string]interface{}{
		"id":   1,
		"name": 2,
		"age":  3,
	}
	q, args := GetSchema("oracle", "select 'colin' name, 1 id from dual where 1=@id and 2=@name", input)
	if !strings.EqualFold(q, "select 'colin' name, 1 id from dual where 1=:1 and 2=:2") {
		t.Error("与期望的值不符:", q)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != 2 {
		t.Error("与期望的值不符:", len(args), args[0], args[1])
	}
	q, args = GetSpSchema("oracle", "grs_p_delivery(@id,@name,@age)", input)
	if !strings.EqualFold(q, "begin grs_p_delivery(:1,:2,:3);end;") {
		t.Error("与期望的值不符:", q)
	}
	if len(args) != 3 || args[0] != 1 || args[1] != 2 || args[2] != 3 {
		t.Error("与期望的值不符:", len(args), args[0], args[1])
	}

}
*/
func TestSchemaExecute(t *testing.T) {
	input := map[string]interface{}{
		"pre_order_no": 1,
	}
	dbMap, err := NewDBMap("oracle", "CY_ESALES/123456@ORCL136")
	if err != nil {
		t.Error(err)
	}
	tr, err := dbMap.Begin()
	if err != nil {
		t.Error(err)
	}
	r, err := tr.Scalar(`select gas_station_name from cy_order_pre where pre_order_no='1'`, input)
	fmt.Println(r.Result)
	if err != nil {
		t.Error(err, r.SQL, len(r.Args))
	}

	buffer, err := json.Marshal(&r.Result)
	if err != nil {
		return
	}
	js := utility.Escape(string(buffer))

	t.Log("json:", js)
	t.Log("r:", r.Result)

	tr.Commit()
}
