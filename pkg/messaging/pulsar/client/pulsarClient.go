package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/apache/pulsar-client-go/pulsar"
)

type consumer struct {
	name          string
	subscribeName string
	topicName     string
	consumerMsg   chan pulsar.ConsumerMessage
	consumer      pulsar.Consumer
}

type Producer struct {
	topicName string
	Producer  pulsar.Producer
}

type pulsarClient struct {
	client    pulsar.Client
	producers []Producer
	consumers []consumer
}

func NewPulsarClient(URL string) (*pulsarClient, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: URL,
	})

	if err != nil {
		return nil, err
	}

	defer client.Close()

	return &pulsarClient{
		client: client,
	}, nil
}

func (pc *pulsarClient) CreateProducer(topic string) error {
	prod, err := pc.client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
		/*Interceptors: pulsar.ProducerInterceptors{
			&tracing.ProducerInterceptor{},
		},*/
	})

	if err != nil {
		return err
	}

	pc.addProducer(topic, prod)

	return nil
}

func (pc *pulsarClient) ExistsGetProducer(topic string) (bool, Producer) {
	for _, p := range pc.producers {
		if p.topicName == topic {
			return true, p
		}
	}
	return false, Producer{}
}

func (pc *pulsarClient) addProducer(topic string, prod pulsar.Producer) {
	exists, _ := pc.ExistsGetProducer(topic)
	if !exists {
		internalProducer := Producer{
			topicName: topic,
			Producer:  prod,
		}
		pc.producers = append(pc.producers, internalProducer)
	}
}

func (pc *pulsarClient) SendMessage(ctx context.Context, topic string, payload interface{}) error {
	exists, producerInternal := pc.ExistsGetProducer(topic)

	if !exists {
		return errors.New("producer for this topic not exists")
	}

	message, err := json.Marshal(payload)

	if err != nil {
		return errors.New("error marshal message")
	}

	_, errSend := producerInternal.Producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: message,
	})
	if errSend != nil {
		return err
	}

	return nil
}

func (pc *pulsarClient) CreateConsumerWithChannels(topic string, subcriptionName string, consumerName string, channelLimit int) error {
	exists, _ := pc.ExistsGetProducer(topic)
	if !exists {
		return errors.New("producer topic not found")
	}

	channel := make(chan pulsar.ConsumerMessage, channelLimit)

	options := pulsar.ConsumerOptions{
		Name:             consumerName,
		Topic:            topic,
		SubscriptionName: subcriptionName,
		Type:             pulsar.Shared,
		/*Interceptors: pulsar.ConsumerInterceptors{
			&tracing.ConsumerInterceptor{},
		},*/
	}

	options.MessageChannel = channel
	cons, err := pc.client.Subscribe(options)
	if err != nil {
		return err
	}

	pc.addConsumer(topic, subcriptionName, consumerName, cons, channel)
	return nil
}

func (pc *pulsarClient) existsConsumer(topic string, subcriptionName string, name string) bool {
	for _, c := range pc.consumers {
		if c.topicName == topic && c.subscribeName == subcriptionName && c.name == name {
			return true
		}
	}
	return false
}

func (pc *pulsarClient) GetConsumer(topic string, subcriptionName string, name string) (pulsar.Consumer, chan pulsar.ConsumerMessage) {
	for _, c := range pc.consumers {
		if c.topicName == topic && c.subscribeName == subcriptionName && c.name == name {
			return c.consumer, c.consumerMsg
		}
	}
	return nil, nil
}

func (pc *pulsarClient) addConsumer(topic string, subcriptionName string, name string, cons pulsar.Consumer, consumerChannel chan pulsar.ConsumerMessage) {
	exists := pc.existsConsumer(topic, subcriptionName, name)
	if !exists {
		internalConsumer := consumer{
			name:          name,
			subscribeName: subcriptionName,
			topicName:     topic,
			consumerMsg:   consumerChannel,
			consumer:      cons,
		}
		pc.consumers = append(pc.consumers, internalConsumer)
	}
}
