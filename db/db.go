package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-oci8"
	_ "github.com/mattn/go-sqlite3"
)

/*
github.com/mattn/go-oci8

http://www.simonzhang.net/?p=2890
http://blog.sina.com.cn/s/blog_48c95a190102w2ln.html
http://www.tudou.com/programs/view/yet9OngrV_4/
https://github.com/wendal/go-oci8/downloads
https://github.com/wendal/go-oci8

安装方法
1. 下载：http://www.oracle.com/technetwork/database/features/instant-client/index.html
2. 解压文件 unzip instantclient-basic-linux.x64-12.1.0.1.0.zip -d /usr/local/
3. 配置环境变量
vi .bash_profile
export ora_home=/usr/local/instantclient_12_1
export PATH=$PATH:$ora_home
export LD_LIBRARY_PATH=$ora_home


*/

const (
	//SQLITE3 Sqlite3数据库
	SQLITE3 = "sqlite3"
	//OCI8 oralce数据库
	OCI8 = "oci8"
)

//DB 数据库实体
type DB struct {
	provider   string
	connString string
	db         *sql.DB
	maxIdle    int
	maxOpen    int
	lang       string
}

//NewDB 创建DB实例
func NewDB(provider string, connString string) (obj *DB, err error) {
	fmt.Println(">-创建DB")
	obj = &DB{provider: provider, connString: connString, lang: "AMERICAN_AMERICA.AL32UTF8"}
	switch strings.ToLower(provider) {
	case "oracle":
		obj.db, err = sql.Open(OCI8, connString)
	case "sqlite":
		obj.db, err = sql.Open(SQLITE3, connString)
	}
	obj.SetPoolSize(0, 0)
	return
}

//SetPoolSize 设置连接池大小
func (db *DB) SetPoolSize(maxIdle int, maxOpen int) {
	if maxIdle != db.maxIdle {
		db.maxIdle = maxIdle
		db.db.SetMaxIdleConns(maxIdle)
	}
	if maxOpen != db.maxOpen {
		db.maxOpen = maxOpen
		db.db.SetMaxOpenConns(maxOpen)
	}
}

//SetLang 设置语言
func (db *DB) SetLang(lang string) {
	db.lang = lang
	db.setEnv("NLS_LANG", lang) //AMERICAN_AMERICA.AL32UTF8
}

//Query 执行SQL查询语句
func (db *DB) Query(query string, args ...interface{}) (dataRows []map[string]interface{}, err error) {
	rows, err := db.db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	dataRows, _, err = queryResolve(rows, 0)
	return

}

//Scalar 执行SQL查询语句
func (db *DB) Scalar(query string, args ...interface{}) (value interface{}, err error) {
	rows, err := db.db.Query(query, args...)
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

func queryResolve(rows *sql.Rows, col int) (dataRows []map[string]interface{}, columns []string, err error) {
	dataRows = make([]map[string]interface{}, 0)
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	for rows.Next() {
		row := make(map[string]interface{})
		dataRows = append(dataRows, row)
		var buffer []interface{}
		for index := 0; index < len(columns); index++ {
			var va []byte
			buffer = append(buffer, &va)
		}
		err = rows.Scan(buffer...)
		if err != nil {
			return
		}
		for index := 0; index < len(columns) && (index < col || col == 0); index++ {
			key := columns[index]
			value := buffer[index]
			if value == nil {
				continue
			} else {
				row[strings.ToLower(key)] = strings.TrimPrefix(fmt.Sprintf("%s", value), "&")
			}
		}
	}
	return
}

//Execute 执行SQL操作语句
func (db *DB) Execute(query string, args ...interface{}) (affectedRow int64, err error) {
	result, err := db.db.Exec(query, args...)
	if err != nil {
		return
	}
	affectedRow, err = result.RowsAffected()
	return
}

//setEnv 设置环境变量
func (db *DB) setEnv(name string, value string) {
	nlsLang := os.Getenv(name)
	if !strings.EqualFold(nlsLang, value) {
		os.Setenv(name, value)
	}
}

//Close 关闭数据库
func (db *DB) Close() {
	db.db.Close()
}
