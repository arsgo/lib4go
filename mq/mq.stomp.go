package mq

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gmallard/stompngo"
)

//StompMQ manage stomp server
type Stomp struct {
	conn    *stompngo.Connection
	cfg     *StompConfig
	Address string
	header  []string
	lk      sync.Mutex
}

//NewStompMQ
func NewStomp(cfg *StompConfig) (mq *Stomp, err error) {
	mq = &Stomp{cfg: cfg}
	mq.header = stompngo.Headers{"accept-version", cfg.AcceptVersion}
	mq.Address = cfg.Address
	err = mq.connect()
	return

}
func (s *Stomp) connect() (err error) {
	con, err := net.Dial("tcp", s.Address)
	if err != nil {
		return
	}
	s.conn, err = stompngo.Connect(con, s.header)
	return
}

//Send
func (s *Stomp) Send(queue string, msg string, timeout int) (err error) {
	s.lk.Lock()
	if !s.conn.Connected() {
		err = s.connect()
	}
	s.lk.Unlock()
	if err != nil {
		return
	}
	header := stompngo.Headers{"destination", queue, "persistent", s.cfg.Persistent}
	if timeout > 0 {
		header = stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", s.cfg.Dest, queue), "persistent", s.cfg.Persistent, "expires",
			fmt.Sprintf("%d000", time.Now().Add(time.Second*time.Duration(timeout)).Unix())}
	}
	return s.conn.Send(header, msg)
}

//Consume
func (s *Stomp) Consume(queue string, call func(MsgHandler)) (err error) {
	s.lk.Lock()
	if !s.conn.Connected() {
		err = s.connect()
	}
	s.lk.Unlock()
	if err != nil {
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
