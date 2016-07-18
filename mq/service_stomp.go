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
	Address       string `json:"address"`
	AcceptVersion string `json:"accept-version"`
	Persistent    string `json:"persistent"`
	Ack           string `json:"ack"`
	Dest          string `json:"dest"`
}

func setDefaultConfig(cfg *StompConfig) (n *StompConfig) {
	n = &StompConfig{}
	n.Address = cfg.Address
	n.AcceptVersion = cfg.AcceptVersion
	n.Persistent = cfg.Persistent
	n.AcceptVersion = cfg.AcceptVersion
	n.Dest = cfg.Dest
	if strings.EqualFold(cfg.AcceptVersion, "") {
		n.AcceptVersion = "1.1"
	}
	if strings.EqualFold(cfg.Persistent, "") {
		n.Persistent = "true"
	}
	if strings.EqualFold(cfg.Ack, "") {
		n.Ack = "client-individual"
	}
	if strings.EqualFold(cfg.Dest, "") {
		n.Dest = "queue"
	}
	return n

}
func NewStompService(sconfig string) (ps IMQService, err error) {
	p := &StompService{}
	ps = p
	err = json.Unmarshal([]byte(sconfig), &p.config)
	if err != nil {
		return
	}
	p.config = setDefaultConfig(p.config)
	if strings.EqualFold(p.config.Address, "") {
		err = errors.New("address is nil")
		return
	}
	p.broker, err = NewStomp(p.config)
	if err != nil {
		return
	}
	return
}
func (k *StompService) Send(queue string, msg string, timeout int) (err error) {
	return k.broker.Send(queue, msg, timeout)
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
