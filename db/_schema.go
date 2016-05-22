package db

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func init() {
	pkg := os.Getenv("PKG_CONFIG_PATH")
	if strings.EqualFold(pkg, "") {
		os.Setenv("PKG_CONFIG_PATH", "/Users/ganyanfei/Projects/golang/pkg-config")
	}
}
func getSchema(format string, data map[string]interface{}, prefix func() string) (query string, args []interface{}) {
	word, _ := regexp.Compile(`@\w+`)
	query = word.ReplaceAllStringFunc(format, func(s string) string {
		args = append(args, data[s])
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
	return getSchema(format, data, f)
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
	case "oci8":
		return GetOracleSchema(format, data)
	case "sqlite3":
		return GetSqliteSchema(format, data)
	}
	return
}
