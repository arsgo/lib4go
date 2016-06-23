package db

import (
	"database/sql"
	"strings"
)

//DBTransaction 事务
type DBTransaction struct {
	tx *sql.Tx
}

//Begin 创建一个事务请求
func (db *DB) Begin() (t *DBTransaction, err error) {
	t = &DBTransaction{}
	t.tx, err = db.db.Begin()
	return
}

//Query 执行查询
func (t *DBTransaction) Query(query string, args ...interface{}) (dataRows []map[string]interface{}, err error) {
	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	dataRows, _, err = queryResolve(rows, 0)
	return
}

//Scalar 执行SQL查询语句
func (t *DBTransaction) Scalar(query string, args ...interface{}) (value interface{}, err error) {
	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	dataRows, columns, err := queryResolve(rows, 1)
	if len(dataRows) > 0 && len(columns) > 0 {
		value = dataRows[0][strings.ToLower(columns[0])]
	}

	return

}

//Execute 执行SQL操作语句
func (t *DBTransaction) Execute(query string, args ...interface{}) (affectedRow int64, err error) {
	result, err := t.tx.Exec(query, args...)
	if err != nil {
		return
	}
	affectedRow, err = result.RowsAffected()
	return
}

//Rollback 回滚所有操作
func (t *DBTransaction) Rollback() error {
	return t.tx.Rollback()
}

//Commit 提交所有操作
func (t *DBTransaction) Commit() error {
	return t.tx.Commit()
}
