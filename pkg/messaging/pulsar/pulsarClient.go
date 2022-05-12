package pulsar

import (
	"context"
	"gounico/pkg/messaging/pulsar/client"

	"github.com/apache/pulsar-client-go/pulsar"
)

type PulsarClient interface {
	CreateProducer(topic string) error
	SendMessage(ctx context.Context, topic string, payload interface{}) error
	CreateConsumerWithChannels(topic string, subcriptionName string, channelLimit int) error
	GetConsumer(topic string, subcriptionName string) (pulsar.Consumer, chan pulsar.ConsumerMessage)
	ExistsGetProducer(topic string) (bool, client.Producer)
}
