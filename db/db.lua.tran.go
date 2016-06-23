package db

import (
	"encoding/json"

	"github.com/yuin/gopher-lua"
)

type DBScriptBindTrans struct {
	db *DBMapTrans
}

func (bind *DBScriptBindTrans) Query(query string, tb *lua.LTable) (r string, err error, sql string, args []interface{}) {
	result, err := bind.db.Query(query, getInputArgs(tb))
	if err != nil {
		return
	}
	sql = result.SQL
	args = result.Args
	buffer, err := json.Marshal(&result.Result)
	if err != nil {
		return
	}
	r = string(buffer)
	return
}

func (bind *DBScriptBindTrans) Scalar(query string, tb *lua.LTable) (r interface{}, err error, sql string, args []interface{}) {
	result, err := bind.db.Scalar(query, getInputArgs(tb))
	if err != nil {
		return
	}
	sql = result.SQL
	args = result.Args
	r = result.Result
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (bind *DBScriptBindTrans) Execute(query string, tb *lua.LTable) (r int64, err error, sql string, args []interface{}) {
	result, err := bind.db.Execute(query, getInputArgs(tb))
	if err != nil {
		return
	}
	r = result.Result
	sql = result.SQL
	args = result.Args
	return
}

//Rollback 回滚所有操作
func (bind *DBScriptBindTrans) Rollback() error {
	return bind.db.Rollback()
}

//Commit 提交所有操作
func (bind *DBScriptBindTrans) Commit() error {
	return bind.db.Commit()
}
