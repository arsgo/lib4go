package mq

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/arsgo/lib4go/concurrent"
	"github.com/gmallard/stompngo"
)

//Stomp stomp
type Stomp struct {
	conn      *stompngo.Connection
	cfg       *StompConfig
	Address   string
	header    []string
	lk        sync.Mutex
	reconnect bool
	mqQueue   *concurrent.ConcurrentMap
}

//NewStomp 构建stomp mq
func NewStomp(cfg *StompConfig) (mq *Stomp, err error) {
	mq = &Stomp{cfg: cfg, Address: cfg.Address}
	mq.header = stompngo.Headers{"accept-version", cfg.AcceptVersion}
	mq.mqQueue = concurrent.NewConcurrentMap()
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

//Send 发送消息
func (s *Stomp) Send(queue string, msg string, timeout int) (err error) {
	index := 0
	reconnect := false
START:
	if index > 3 {
		return
	}
	s.lk.Lock()
	if reconnect || !s.conn.Connected() {
		s.Close()
		err = s.connect()
	}
	s.lk.Unlock()
	if err != nil {
		reconnect = true
		goto START
	}
	header := stompngo.Headers{"destination", queue, "persistent", s.cfg.Persistent}
	if timeout > 0 {
		header = stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", s.cfg.Dest, queue), "persistent", s.cfg.Persistent, "expires",
			fmt.Sprintf("%d000", time.Now().Add(time.Second*time.Duration(timeout)).Unix())}
	}
	err = s.conn.Send(header, msg)
	if err != nil {
		reconnect = true
		index++
		goto START
	}
	return
}
func (s *Stomp) createConsumer(p ...interface{}) (ch interface{}, err error) {
	queue := p[0].(string)
	subscriberHeader := stompngo.Headers{"destination", fmt.Sprintf("/%s/%s", s.cfg.Dest, queue), "ack", s.cfg.Ack}
	msgChan, err := s.conn.Subscribe(subscriberHeader)
	if err != nil {
		return
	}
	ch = msgChan
	return
}

//Consume 注册消费信息
func (s *Stomp) Consume(queue string, call func(MsgHandler)) (err error) {
	s.lk.Lock()
	if !s.conn.Connected() {
		err = s.connect()
	}
	s.lk.Unlock()
	if err != nil {
		return
	}
	success, ch, err := s.mqQueue.GetOrAdd(queue, s.createConsumer, queue)
	if err != nil {
		return
	}
	if !success {
		err = fmt.Errorf("重复订阅消息:%s", queue)
		return
	}
	msgChan := ch.(<-chan stompngo.MessageData)
START:
	for {
		select {
		case msg, ok := <-msgChan:
			if ok {
				call(NewStompMessage(s, &msg.Message))
			} else {
				break START
			}
		}
	}
	return
}

//UnConsume 取消注册消费
func (s *Stomp) UnConsume(queue string) {
	subscriberHeader := stompngo.Headers{"destination",
		fmt.Sprintf("/%s/%s", s.cfg.Dest, queue), "ack", s.cfg.Ack}
	s.conn.Unsubscribe(subscriberHeader)
	_, ok := s.mqQueue.Get(queue)
	if !ok {
		return
	}

}

//Close 关闭当前连接
func (s *Stomp) Close() {
	if !s.conn.Connected() {
		return
	}
	s.conn.Disconnect(stompngo.Headers{})
}
