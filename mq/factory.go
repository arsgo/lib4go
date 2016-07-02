package mq

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	stompMQ = "stomp"
	kafkaMQ = "kafka"
)

type MsgHandler interface {
	Ack()
	Nack()
	GetMessage() string
}

type IMQService interface {
	Consume(string, func(MsgHandler)) error
	UnConsume(string)
	Send(string, string, int) error
	Close()
}

type MQConfig struct {
	Type string `json:"type"`
}

func NewMQService(config string) (svs IMQService, err error) {
	p := &MQConfig{}
	err = json.Unmarshal([]byte(config), &p)
	if err != nil {
		return
	}
	switch strings.ToLower(p.Type) {
	case stompMQ:
		svs, err = NewStompService(config)
	//case kafkaMQ:
	//	svs, err = NewKafkaService(config)
	default:
		err = errors.New("not support mq type")
	}
	return

}
