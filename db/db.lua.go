package db

import (
	"encoding/json"

	"github.com/yuin/gopher-lua"
)

type DBScriptBind struct {
	db *DBMap
}

func NewDBScriptBind(config string) (b *DBScriptBind, err error) {
	b = &DBScriptBind{}
	b.db, err = NewDBMapByConfig(config)
	if err != nil {
		return
	}
	return
}

func getInputArgs(tb *lua.LTable) (data map[string]interface{}) {
	data = make(map[string]interface{})
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		data[key.String()] = value.String()
	})
	return
}

//Query 根据包含@名称占位符的查询语句执行查询语句
func (bind *DBScriptBind) Query(query string, tb *lua.LTable) (r string, err error, sql string, args []interface{}) {
	result, err := bind.db.Query(query, getInputArgs(tb))
	sql = result.SQL
	args = result.Args
	if result.Result == nil {
		return
	}
	buffer, err := json.Marshal(&result.Result)
	if err != nil {
		return
	}
	r = string(buffer)
	return
}

//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (bind *DBScriptBind) Scalar(query string, tb *lua.LTable) (r interface{}, err error, sql string, args []interface{}) {
	result, err := bind.db.Scalar(query, getInputArgs(tb))
	r = result.Result
	sql = bind.db.GetReplaceSchema(result.SQL, result.Args)
	args = result.Args
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (bind *DBScriptBind) Execute(query string, tb *lua.LTable) (r int64, err error, sql string, args []interface{}) {
	result, err := bind.db.Execute(query, getInputArgs(tb))
	r = result.Result
	sql = bind.db.GetReplaceSchema(result.SQL, result.Args)
	args = result.Args
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (bind *DBScriptBind) ExecuteSP(query string, tb *lua.LTable) (r int64, err error, sql string, args []interface{}) {
	result, err := bind.db.ExecuteSP(query, getInputArgs(tb))
	r = result.Result
	sql = bind.db.GetReplaceSchema(result.SQL, result.Args)
	args = result.Args
	return
}

//Begin 开始一个事务
func (bind *DBScriptBind) Begin() (bts *DBScriptBindTrans, err error) {
	ts, err := bind.db.Begin()
	if err != nil {
		return
	}
	bts = &DBScriptBindTrans{db: ts}
	return
}
