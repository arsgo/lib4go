package mq

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gmallard/stompngo"
)

//StompMQ manage stomp server
type Stomp struct {
	conn *stompngo.Connection
	cfg  *StompConfig
}

//NewStompMQ
func NewStomp(cfg *StompConfig) (mq *Stomp, err error) {
	mq = &Stomp{cfg: cfg}
	con, err := net.Dial("tcp", cfg.Address)
	if err != nil {
		return
	}
	header := stompngo.Headers{"accept-version", cfg.AcceptVersion}
	mq.conn, err = stompngo.Connect(con, header)
	return
}

//Send
func (s *Stomp) Send(queue string, msg string, timeout int) error {
	header := stompngo.Headers{"destination", queue, "persistent", s.cfg.Persistent}
	if timeout > 0 {
		header = stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", s.cfg.Dest, queue), "persistent", s.cfg.Persistent, "expires",
			fmt.Sprintf("%d000", time.Now().Add(time.Second*time.Duration(timeout)).Unix())}
	}
	return s.conn.Send(header, msg)
}

//Subscribe
func (s *Stomp) Consume(queue string, call func(MsgHandler)) (err error) {
	if !s.conn.Connected() {
		err = errors.New("not connect to stomp server")
		return
	}
	subscriberHeader := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", s.cfg.Dest, queue), "ack", s.cfg.Ack}
	msgChan, err := s.conn.Subscribe(subscriberHeader)
	if err != nil {
		return
	}

	for {
		msg := <-msgChan
		call(NewStompMessage(s, &msg.Message))
	}
}
func (s *Stomp) UnConsume(queue string) {
	subscriberHeader := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", s.cfg.Dest, queue), "ack", s.cfg.Ack}
	s.conn.Unsubscribe(subscriberHeader)
}

//Close
func (s *Stomp) Close() {
	if !s.conn.Connected() {
		return
	}
	s.conn.Disconnect(stompngo.Headers{})
}
