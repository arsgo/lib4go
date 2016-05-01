package influxdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/colinyl/lib4go/utility"
)

type influxDbConfig struct {
	Address   string `json:"address"`
	DbName    string `json:"db"`
	UserName  string `json:"user"`
	Password  string `json:"password"`
	RowFormat string `json:"row"`
}

func Save(f string, rows []map[string]interface{}) (err error) {
	config := &influxDbConfig{}
	err = json.Unmarshal([]byte(f), &config)
	if err != nil {
		return
	}
	if strings.EqualFold(config.Address, "") ||
		strings.EqualFold(config.DbName, "") ||
		strings.EqualFold(config.RowFormat, "") {
		err = errors.New("influxDbConfig必须参数不能为空")
		return
	}
	url := fmt.Sprintf("%s/write?db=%s", config.Address, config.DbName)
	var datas []string
	for i := 0; i < len(rows); i++ {
		d := utility.NewDataMaps(rows[i])
		datas = append(datas, d.Translate(config.RowFormat))
	}
	data := strings.Join(datas, "\n")
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil
	}
	err = errors.New(fmt.Sprintf("error:%d", resp.StatusCode))
	return
}
