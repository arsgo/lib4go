package utility

import (
	"encoding/json"
	"net/url"
)

//GetParams 输入参数为URL参数，转换为Json字符串
func GetParams(urlQuery string) (res string, err error) {
	values, err := url.ParseQuery(urlQuery)
	if err != nil {
		return
	}
	result := make(map[string]interface{})
	for k, v := range values {
		if len(v) == 1 {
			result[k] = v[0]
		} else {
			result[k] = v
		}
	}

	buffer, err := json.Marshal(&result)
	if err != nil {
		return
	}
	return string(buffer), nil
}
