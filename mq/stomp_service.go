package mq

import (
	"encoding/json"
	"errors"
	"strings"
)

type StompService struct {
	config *StompConfig
	broker *Stomp
}

type StompConfig struct {
	Address string `json:"address"`
}

func NewStompService(sconfig string) (ps IMQService, err error) {
	p := &StompService{}
	ps = p
	err = json.Unmarshal([]byte(sconfig), &p.config)
	if err != nil {
		return
	}
	if strings.EqualFold(p.config.Address, "") {
		err = errors.New("address is nil")
		return
	}
	p.broker, err = NewStomp(p.config.Address)
	if err != nil {
		return
	}
	return
}
func (k *StompService) Send(queue string, msg string) (err error) {
	return k.broker.Send(queue, msg)
}

func (k *StompService) Consume(queue string, callback func(MsgHandler)) (err error) {
	return k.broker.Consume(queue, callback)
}
func (k *StompService) UnConsume(queue string) {
	k.broker.UnConsume(queue)
}
func (k *StompService) Close() {
	k.broker.Close()
}
