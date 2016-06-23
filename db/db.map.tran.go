package db

//DBMapTrans 
type DBMapTrans struct {
	db       *DBTransaction
	provider string
}

//Begin 创建事务
func (db *DBMap) Begin() (t *DBMapTrans, err error) {
	t = &DBMapTrans{}
	t.provider = db.db.provider
	t.db, err = db.db.Begin()
	return
}

//Query 根据包含@名称占位符的查询语句执行查询语句
func (db *DBMapTrans) Query(query string, data map[string]interface{}) (r DBQueryResult, err error) {
	r.SQL, r.Args = GetSchema(db.provider, query, data)
	r.Result, err = db.db.Query(r.SQL, r.Args...)
	return
}
//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (db *DBMapTrans) Scalar(query string, data map[string]interface{}) (r DBScalarResult, err error) {
	r.SQL, r.Args = GetSchema(db.provider, query, data)
	r.Result, err = db.db.Scalar(r.SQL, r.Args...)
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (db *DBMapTrans) Execute(query string, data map[string]interface{}) (r DBExecuteResult, err error) {
	r.SQL, r.Args = GetSchema(db.provider, query, data)
	r.Result, err = db.db.Execute(r.SQL, r.Args...)
	return
}

//Rollback 回滚所有操作
func (db *DBMapTrans) Rollback() error {
	return db.db.Rollback()
}

//Commit 提交所有操作
func (db *DBMapTrans) Commit() error {
	return db.db.Commit()
}
