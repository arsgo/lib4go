package db

//DBMAP
type DBMap struct {
	db *DB
}

//DBQueryResult
type DBQueryResult struct {
	SQL    string
	Args   []interface{}
	Result []map[string]interface{}
}

//DBExecuteResult
type DBExecuteResult struct {
	SQL    string
	Args   []interface{}
	Result int64
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
