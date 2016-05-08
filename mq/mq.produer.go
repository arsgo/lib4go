package mq

type MQProducer struct {
	service IMQService
}

func (p *MQProducer) Send(queue string, content string) error {
	return p.service.Send(queue, content)
}

func NewMQProducer(config string) (m *MQProducer, err error) {
	m = &MQProducer{}
	m.service, err = NewMQService(config)
	return
}
