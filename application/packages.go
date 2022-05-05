package application

import (
	"gounico/config"
	"gounico/infrastructure/database/dynamodb"
	"gounico/infrastructure/database/dynamodb/client"
	"gounico/infrastructure/logging"
	"gounico/infrastructure/logging/zap"
	"gounico/infrastructure/messaging/pulsar"
	clientPulsar "gounico/infrastructure/messaging/pulsar/client"

	"go.uber.org/fx"
)

var PackagesModule = fx.Provide(
	NewLogger,
	NewDynamoClient,
	NewPulsarClient,
)

func NewLogger(config config.Configuration) (logging.Logger, error) {
	logger, err := zap.NewZapLogger(config.Server.LogPath)
	return logger, err
}

func NewDynamoClient(config config.Configuration) dynamodb.DynamoClient {
	clientDynamo := client.NewDynamoDBClient(config.Database.EndpointURL,
		config.Database.Region,
		config.Database.AccessKeyID,
		config.Database.SecretAccessKey,
		config.Database.SessionToken,
		config.Database.Maintable)
	return clientDynamo
}

func NewPulsarClient(config config.Configuration) (pulsar.PulsarClient, error) {
	pulsarClient, err := clientPulsar.NewPulsarClient(config.Messaging.BrokerURL)

	if err != nil {
		return nil, err
	}

	for _, configConsumer := range config.Messaging.Configurations {
		errProd := pulsarClient.CreateProducer(configConsumer.Topic)
		if errProd != nil {
			return nil, errProd
		}

		errCons := pulsarClient.CreateConsumerWithChannels(configConsumer.Topic, configConsumer.Subscriber, config.Messaging.ConsumerLimit)
		if errCons != nil {
			return nil, errCons
		}

	}

	return pulsarClient, err
}
