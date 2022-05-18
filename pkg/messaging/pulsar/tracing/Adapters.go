package tracing

import (
	"github.com/apache/pulsar-client-go/pulsar"
)

type ProducerMessageAdapter struct {
	Message *pulsar.ProducerMessage
}

func (a *ProducerMessageAdapter) Get(key string) string {
	return a.Message.Properties[key]
}

func (a *ProducerMessageAdapter) Set(key string, value string) {
	a.Message.Properties[key] = value
}

func (a *ProducerMessageAdapter) Keys() []string {
	var keys []string
	for key, _ := range a.Message.Properties {
		keys = append(keys, key)
	}
	return keys
}

type ConsumerMessageAdapter struct {
	Message pulsar.ConsumerMessage
}

func (a *ConsumerMessageAdapter) Get(key string) string {
	return a.Message.Properties()[key]
}

func (a *ConsumerMessageAdapter) Set(key string, value string) {

}

func (a *ConsumerMessageAdapter) Keys() []string {
	var keys []string
	for key, _ := range a.Message.Properties() {
		keys = append(keys, key)
	}
	return keys
}
