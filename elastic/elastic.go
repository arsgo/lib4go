package elastic

import (
	"encoding/json"
	"errors"

	"github.com/mattbaird/elastigo/lib"
)

type elasticHost struct {
	host []string
	conn *elastigo.Conn
}

func NewElastic(hosts []string) (host *elasticHost) {
	host = &elasticHost{}
	host.conn = elastigo.NewConn()
	host.host = hosts
	host.conn.SetHosts(hosts)
	return host
}

func (host *elasticHost) Create(name string, typeName string, jsonData string) (err error) {
	response, err := host.conn.Index(name, typeName, "0", nil, jsonData)
	if err != nil {
		return
	}
	host.conn.Flush()
	if response.Ok {
		return
	}
	return errors.New("")
}
func (host *elasticHost) Search(name string, typeName string, query string) (result string, err error) {
	out, err := host.conn.Search(name, typeName, nil, query)
	if err != nil {
		return
	}
	var resultLst []*json.RawMessage
	for i := 0; i < len(out.Hits.Hits); i++ {
		resultLst = append(resultLst, (out.Hits.Hits[i].Source))
	}
	buffer, err := json.Marshal(&resultLst)
	if err != nil {
		return
	}
	result = string(buffer)
	return
}
