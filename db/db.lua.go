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

func (bind *DBScriptBind) getInputArgs(tb *lua.LTable) (data map[string]interface{}) {
	data = make(map[string]interface{})
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		data[key.String()] = value.String()
	})
	return
}

//Query 根据包含@名称占位符的查询语句执行查询语句
func (bind *DBScriptBind) Query(query string, tb *lua.LTable) (r string, err error) {
	result, err := bind.db.Query(query, bind.getInputArgs(tb))
	buffer, err := json.Marshal(&result)
	if err != nil {
		return
	}
	r = string(buffer)
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (bind *DBScriptBind) Execute(query string, tb *lua.LTable) (r string, err error) {
	result, err := bind.db.Execute(query, bind.getInputArgs(tb))
	buffer, err := json.Marshal(&result)
	if err != nil {
		return
	}
	r = string(buffer)
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (bind *DBScriptBind) ExecuteSP(query string, tb *lua.LTable) (r string, err error) {
	result, err := bind.db.ExecuteSP(query, bind.getInputArgs(tb))
	buffer, err := json.Marshal(&result)
	if err != nil {
		return
	}
	r = string(buffer)
	return
}
