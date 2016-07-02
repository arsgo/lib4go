package mq

type MQProducer struct {
	service IMQService
}

func (p *MQProducer) Send(queue string, content string, timeout int) error {
	return p.service.Send(queue, content, timeout)
}

func NewMQProducer(config string) (m *MQProducer, err error) {
	m = &MQProducer{}
	m.service, err = NewMQService(config)
	return
}
