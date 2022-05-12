package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/apache/pulsar-client-go/pulsar"
)

type consumer struct {
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

func (pc *pulsarClient) CreateConsumerWithChannels(topic string, subcriptionName string, channelLimit int) error {
	exists, _ := pc.ExistsGetProducer(topic)
	if !exists {
		return errors.New("producer topic not found")
	}

	channel := make(chan pulsar.ConsumerMessage, channelLimit)

	options := pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: subcriptionName,
		Type:             pulsar.Exclusive,
	}

	options.MessageChannel = channel
	cons, err := pc.client.Subscribe(options)
	if err != nil {
		return err
	}

	pc.addConsumer(topic, subcriptionName, cons, channel)
	return nil
}

func (pc *pulsarClient) existsConsumer(topic string, subcriptionName string) bool {
	for _, c := range pc.consumers {
		if c.topicName == topic && c.subscribeName == subcriptionName {
			return true
		}
	}
	return false
}

func (pc *pulsarClient) GetConsumer(topic string, subcriptionName string) (pulsar.Consumer, chan pulsar.ConsumerMessage) {
	for _, c := range pc.consumers {
		if c.topicName == topic && c.subscribeName == subcriptionName {
			return c.consumer, c.consumerMsg
		}
	}
	return nil, nil
}

func (pc *pulsarClient) addConsumer(topic string, subcriptionName string, cons pulsar.Consumer, consumerChannel chan pulsar.ConsumerMessage) {
	exists := pc.existsConsumer(topic, subcriptionName)
	if !exists {
		internalConsumer := consumer{
			subscribeName: subcriptionName,
			topicName:     topic,
			consumerMsg:   consumerChannel,
			consumer:      cons,
		}

		pc.consumers = append(pc.consumers, internalConsumer)
	}
}
