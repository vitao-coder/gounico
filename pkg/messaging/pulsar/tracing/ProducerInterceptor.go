package tracing

import (
	"github.com/apache/pulsar-client-go/pulsar"
)

type ProducerInterceptor struct {
}

func (t *ProducerInterceptor) BeforeSend(producer pulsar.Producer, message *pulsar.ProducerMessage) {
}

func (t *ProducerInterceptor) OnSendAcknowledgement(producer pulsar.Producer, message *pulsar.ProducerMessage, msgID pulsar.MessageID) {
}
