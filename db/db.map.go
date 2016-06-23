package db

import (
	"encoding/json"
	"strings"
)

//DBMAP
type DBMap struct {
	db *DB
}
type DBMapConfig struct {
	Provider   string `json:"provider"`
	ConnString string `json:"connString"`
	Min        int    `json:"min"`
	Max        int    `json:"max"`
	Lang       string `json:"lang"`
}

//DBQueryResult
type DBQueryResult struct {
	SQL    string                   `json:"sql"`
	Args   []interface{}            `json:"args "`
	Result []map[string]interface{} `json:"data"`
}

//DBQueryResult
type DBScalarResult struct {
	SQL    string        `json:"sql"`
	Args   []interface{} `json:"args "`
	Result interface{}   `json:"data"`
}

//DBExecuteResult
type DBExecuteResult struct {
	SQL    string        `json:"sql"`
	Args   []interface{} `json:"args "`
	Result int64         `json:"result"`
}

//NewDBMapByConfig 根据json配置文件创建DBMAP
func NewDBMapByConfig(config string) (obj *DBMap, err error) {
	var dbConfig DBMapConfig
	err = json.Unmarshal([]byte(config), &dbConfig)
	if err != nil {
		return
	}
	obj, err = NewDBMap(dbConfig.Provider, dbConfig.ConnString)
	if err != nil {
		return
	}
	if dbConfig.Min > 0 && dbConfig.Max > 0 {
		obj.SetPoolSize(dbConfig.Min, dbConfig.Max)
	}
	if !strings.EqualFold(dbConfig.Lang, "") {
		obj.SetLang(dbConfig.Lang)
	}
	return
}

//NewDBMap 构建DBMAP
func NewDBMap(provider string, connString string) (obj *DBMap, err error) {
	obj = &DBMap{}
	obj.db, err = NewDB(provider, connString)
	return
}

//SetPoolSize 设置连接池大小
func (db *DBMap) SetPoolSize(maxIdle int, maxOpen int) {
	db.db.SetPoolSize(maxIdle, maxOpen)
}

//SetLang 设置语言
func (db *DBMap) SetLang(lang string) {
	db.SetLang(lang)
}

//Query 根据包含@名称占位符的查询语句执行查询语句
func (db *DBMap) Query(query string, data map[string]interface{}) (r DBQueryResult, err error) {
	r.SQL, r.Args = GetSchema(db.db.provider, query, data)
	r.Result, err = db.db.Query(r.SQL, r.Args...)
	return
}

//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (db *DBMap) Scalar(query string, data map[string]interface{}) (r DBScalarResult, err error) {
	r.SQL, r.Args = GetSchema(db.db.provider, query, data)
	r.Result, err = db.db.Scalar(r.SQL, r.Args...)
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (db *DBMap) Execute(query string, data map[string]interface{}) (r DBExecuteResult, err error) {
	r.SQL, r.Args = GetSchema(db.db.provider, query, data)
	r.Result, err = db.db.Execute(r.SQL, r.Args...)
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (db *DBMap) ExecuteSP(query string, data map[string]interface{}) (r DBExecuteResult, err error) {
	r.SQL, r.Args = GetSpSchema(db.db.provider, query, data)
	r.Result, err = db.db.Execute(r.SQL, r.Args...)
	return
}
