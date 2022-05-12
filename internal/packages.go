package internal

import (
	"gounico/config"
	"gounico/pkg/database/dynamodb"
	"gounico/pkg/database/dynamodb/client"
	"gounico/pkg/logging"
	"gounico/pkg/logging/zap"
	"gounico/pkg/messaging/pulsar"
	clientPulsar "gounico/pkg/messaging/pulsar/client"
	"net/http"
	"time"

	"go.uber.org/fx"
)

var PackagesModule = fx.Provide(
	NewLogger,
	NewDynamoClient,
	NewPulsarClient,
	NewHttpClient,
)

func NewLogger() (logging.Logger, error) {
	logger, err := zap.NewZapLogger()
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

	}

	return pulsarClient, err
}

func NewHttpClient(config config.Configuration) *http.Client {
	clientHttp := &http.Client{
		Timeout: 30 * time.Second,
	}
	return clientHttp
}
