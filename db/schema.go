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
		value := data[s[1:]]
		if value != nil && !strings.EqualFold(fmt.Sprintf("%s", value), "") {
			args = append(args, value)
		} else {
			args = append(args, nil)
		}
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
func ReplaceOracleSchema(query string, args []interface{}) (r string) {	
	r = query
	if strings.EqualFold(query, "") || args == nil || len(args) == 0 {
		return
	}
	for i, v := range args {
		r = strings.Replace(r, fmt.Sprintf(":%d", i+1), fmt.Sprintf("'%s'", v), -1)
	}
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
func GetReplaceSpSchema(provider string, query string, args []interface{}) (q string) {
	p := strings.ToLower(provider)
	switch p {
	case "oracle":
		return ReplaceOracleSchema(query, args)
	case "sqlite3":
		return ""
	}
	return
}
