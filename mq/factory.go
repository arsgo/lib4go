package mq

import (
	"encoding/json"
	"errors"
)

const (
	stompMQ = "stomp"
)

type MQConfig struct {
	Type string `json:"type"`
}

func NewMQService(config string) (svs IMQService, err error) {
	p := &MQConfig{}
	err = json.Unmarshal([]byte(config), &p)
	if err != nil {
		return
	}
	switch p.Type {
	case stompMQ:
		svs, err = NewStompService(config)
		return
	}
	err = errors.New("not support mq type")
	return

}
