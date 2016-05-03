package mq

import s "github.com/gmallard/stompngo"

type StompMessage struct {
	s       *Stomp
	msg     *s.Message
	Message string
}

//Ack
func (m *StompMessage) Ack() {
	m.s.conn.Ack(m.msg.Headers)
}
func (m *StompMessage) Nack() {
	m.s.conn.Nack(m.msg.Headers)
}
func (m *StompMessage) GetMessage() string {
	return m.Message
}

//NewMessage
func NewMessage(s *Stomp, msg *s.Message) *StompMessage {
	return &StompMessage{s: s, msg: msg, Message: string(msg.Body)}
}
