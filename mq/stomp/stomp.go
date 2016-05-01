package stomp

import (
	"net"

	//	st "github.com/gmallard/stompngo"

	"github.com/gmallard/stompngo"
	//"github.com/go-stomp/stomp"

	//"github.com/gmallard/stompngo"
)

//StompMQ manage stomp server
type Stomp struct {
	conn    *stompngo.Connection
	address string
}

//NewStompMQ
func NewStomp(address string) (mq *Stomp, err error) {
	mq = &Stomp{address: address}
	con, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	header := stompngo.Headers{"accept-version", "1.1"}
	mq.conn, err = stompngo.Connect(con, header)
	return
}

//Send
func (s *Stomp) Send(queue string, msg string) error {
	header := stompngo.Headers{"destination", queue, "persistent", "true"}
	return s.conn.Send(header, msg)
}

/*
//Subscribe
func (s *Stomp) Consume(queue string, call func(MsgHandler)) (err error) {
	if !s.conn.Connected() {
		err = errors.New("not connect to stomp server")
		return
	}
	subscriberHeader := stompngo.Headers{"destination",
		fmt.Sprintf("/queue/%s", queue), "ack", "client"}
	msgChan, err := s.conn.Subscribe(subscriberHeader)
	if err != nil {
		return
	}

	for {
		<-msgChan
		//call(NewMessage(s, msg))
	}
}
*/

//Close
func (s *Stomp) Close() {
	s.conn.Disconnect(stompngo.Headers{})
}
