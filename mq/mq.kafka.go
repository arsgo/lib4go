package mq

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jdamick/kafka"
)

type Kafka struct {
	producerAddress string
	consumerAddress string
	topic           string
	publisher       *kafka.BrokerPublisher
	consumer        *kafka.BrokerConsumer
	consumeQueue    chan *kafka.Message
	concurrent      int
	quitChan        chan struct{}
}

func NewKafka(producerAddress string, consumerAddress string, topic string, partition int, concurrent int) (mq *Kafka, err error) {
	mq = &Kafka{producerAddress: producerAddress, consumerAddress: consumerAddress,
		topic: topic, consumeQueue: make(chan *kafka.Message, concurrent),
		quitChan: make(chan struct{}, 0), concurrent: concurrent}
	mq.publisher = kafka.NewBrokerPublisher(producerAddress, topic, partition)
	mq.consumer = kafka.NewBrokerConsumer(consumerAddress, topic, 0, 0, 1048576)
	return
}

func (k *Kafka) Send(topic string, content string) (err error) {
	if !strings.EqualFold(topic, k.topic) {
		err = errors.New(fmt.Sprintf("topic not support:%s", topic))
		return
	}
	_, err = k.publisher.Publish(kafka.NewMessage([]byte(content)))
	return
}

//Subscribe
func (k *Kafka) Consume(queue string, call func(MsgHandler)) (err error) {
	go k.consumer.ConsumeOnChannel(k.consumeQueue, 10, k.quitChan)
LOOP:
	for {
		select {
		case <-k.quitChan:
			break LOOP
		case msg := <-k.consumeQueue:
			call(NewKafkaMessage(msg))
		}
	}
	return nil
}
func (k *Kafka) UnConsume(queue string) {
	if !strings.EqualFold(queue, k.topic) {
		return
	}
	k.quitChan <- struct{}{}
}

//Close
func (k *Kafka) Close() {
	k.quitChan <- struct{}{}
}
