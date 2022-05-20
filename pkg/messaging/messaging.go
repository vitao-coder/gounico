package messaging

import (
	"context"
	"gounico/pkg/messaging/pulsar/client"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Messaging interface {
	CreateProducer(topic string) error
	SendMessage(ctx context.Context, topic string, payload interface{}) error
	CreateConsumerWithChannels(topic string, subcriptionName string, consumerName string, channelLimit int) error
	GetConsumer(topic string, subcriptionName string, name string) (pulsar.Consumer, chan pulsar.ConsumerMessage)
	ExistsGetProducer(topic string) (bool, client.Producer)
}
