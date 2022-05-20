package internal

import (
	"gounico/config"
	"gounico/global"
	"gounico/pkg/database"
	"gounico/pkg/database/dynamodb/client"
	"gounico/pkg/logging"
	"gounico/pkg/logging/zap"
	"gounico/pkg/messaging"
	clientPulsar "gounico/pkg/messaging/pulsar/client"
	"gounico/pkg/telemetry"
	"gounico/pkg/telemetry/openTelemetry"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"go.uber.org/fx"
)

var PackagesModule = fx.Provide(
	NewLogger,
	NewDynamoClient,
	NewPulsarClient,
	NewHttpClient,
	NewOpenTelemetry,
)

func NewLogger() (logging.Logger, error) {
	logger, err := zap.NewZapLogger()
	return logger, err
}

func NewDynamoClient(config config.Configuration) database.Database {
	clientDynamo := client.NewDynamoDBClient(config.Database.EndpointURL,
		config.Database.Region,
		config.Database.AccessKeyID,
		config.Database.SecretAccessKey,
		config.Database.SessionToken,
		config.Database.Maintable)
	return clientDynamo
}

func NewPulsarClient(config config.Configuration) (messaging.Messaging, error) {
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

func NewHttpClient() *http.Client {
	clientHttp := &http.Client{
		Timeout:   60 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	return clientHttp
}

func NewOpenTelemetry(config config.Configuration) telemetry.OpenTelemetry {
	global.AppName = config.Telemetry.AppName
	tracerProvider := openTelemetry.NewTracer(config.Telemetry.JaegerEndpoint, config.Telemetry.AppName)
	return tracerProvider
}
