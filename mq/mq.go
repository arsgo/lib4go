package mq

import "github.com/colinyl/lib4go/mq/stomp"

type IMQService interface {
	Consume(string, func(stomp.MsgHandler) bool) error
	Send(string, string) error
	Close()
}
