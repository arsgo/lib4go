package mq

import (
	"fmt"

	s "github.com/gmallard/stompngo"
)

type StompMessage struct {
	s       *Stomp
	msg     *s.Message
	Message string
}

//Ack
func (m *StompMessage) Ack() {
	fmt.Println("------ack:", m.msg.Headers)
	m.s.conn.Ack(m.msg.Headers)
}
func (m *StompMessage) Nack() {
	//	fmt.Println("------nack:", m.msg.Headers)
	//	m.s.conn.Nack(m.msg.Headers)
}
func (m *StompMessage) GetMessage() string {
	return m.Message
}

//NewMessage
func NewStompMessage(s *Stomp, msg *s.Message) *StompMessage {
	return &StompMessage{s: s, msg: msg, Message: string(msg.Body)}
}
