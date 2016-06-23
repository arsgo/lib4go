package db

import (
	"fmt"
	"regexp"
	"strings"
)

func getSchema(format string, data map[string]interface{}, prefix func() string) (query string, args []interface{}) {
	args = make([]interface{}, 0)
	word, _ := regexp.Compile(`@\w+`)
	query = word.ReplaceAllStringFunc(format, func(s string) string {
		args = append(args, data[s[1:]])
		return prefix()
	})
	return
}

//GetOracleSchema 获取ORACLE参数结构
func GetOracleSchema(format string, data map[string]interface{}) (query string, args []interface{}) {
	index := 0
	f := func() string {
		index++
		return fmt.Sprint(":", index)
	}
	query, args = getSchema(format, data, f)
	return

}

//GetOracleSPSchema 获取ORACLE参数结构
func GetOracleSPSchema(format string, data map[string]interface{}) (query string, args []interface{}) {
	index := 0
	f := func() string {
		index++
		return fmt.Sprint(":", index)
	}
	query, args = getSchema(format, data, f)
	query = fmt.Sprintf("begin %s;end;", query)
	return

}

//GetSqliteSchema 获取Sqlite参数结构
func GetSqliteSchema(format string, data map[string]interface{}) (query string, args []interface{}) {
	f := func() string {
		return "?"
	}
	return getSchema(format, data, f)
}

//GetSchema 获取数据库结构
func GetSchema(provider string, format string, data map[string]interface{}) (query string, args []interface{}) {
	p := strings.ToLower(provider)
	switch p {
	case "oracle":
		return GetOracleSchema(format, data)
	case "sqlite3":
		return GetSqliteSchema(format, data)
	}
	return
}

//GetSpSchema 获取存储过程结构
func GetSpSchema(provider string, format string, data map[string]interface{}) (query string, args []interface{}) {
	p := strings.ToLower(provider)
	switch p {
	case "oracle":
		return GetOracleSPSchema(format, data)
	case "sqlite3":
		return GetSqliteSchema(format, data)
	}
	return
}
