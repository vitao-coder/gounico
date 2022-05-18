package tracing

import (
	"github.com/apache/pulsar-client-go/pulsar"
)

type ConsumerInterceptor struct {
}

func (t *ConsumerInterceptor) BeforeConsume(message pulsar.ConsumerMessage) {
}

func (t *ConsumerInterceptor) OnAcknowledge(consumer pulsar.Consumer, msgID pulsar.MessageID) {}

func (t *ConsumerInterceptor) OnNegativeAcksSend(consumer pulsar.Consumer, msgIDs []pulsar.MessageID) {
}
