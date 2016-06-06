package mq
/*
import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type KafkaService struct {
	config *kafkaConfig
	broker *Kafka
}

type kafkaConfig struct {
	ProducerAddress string `json:"producer"`
	ConsumerAddress string `json:"consumer"`
	Topic           string `json:"topic"`
	Partition       int    `json:"partition"`
	Concurrent      int    `json:"concurrent"`
}

func NewKafkaService(sconfig string) (ps IMQService, err error) {
	p := &KafkaService{}
	ps = p
	err = json.Unmarshal([]byte(sconfig), &p.config)
	if err != nil {
		return
	}
	if strings.EqualFold(p.config.ProducerAddress, "") ||
		strings.EqualFold(p.config.ConsumerAddress, "") ||
		strings.EqualFold(p.config.Topic, "") {
		err = errors.New(fmt.Sprint("producer or consumer  or topic is nil in:", sconfig))
		return
	}
	p.broker, err = NewKafka(p.config.ProducerAddress, p.config.ConsumerAddress,
		p.config.Topic, p.config.Partition, p.config.Concurrent)
	if err != nil {
		return
	}
	return
}
func (k *KafkaService) Send(queue string, msg string) (err error) {
	return k.broker.Send(queue, msg)
}

func (k *KafkaService) Consume(queue string, callback func(MsgHandler)) (err error) {
	return k.broker.Consume(queue, callback)
}
func (k *KafkaService) UnConsume(queue string) {
	k.broker.UnConsume(queue)
}
func (k *KafkaService) Close() {
	k.broker.Close()
}
*/