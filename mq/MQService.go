package mq

type MsgHandler interface {
	Ack()
	Nack()
	GetMessage() string
}

type IMQService interface {
	Consume(string, func(MsgHandler)) error
	UnConsume(string)
	Send(string, string) error
	Close()
}
