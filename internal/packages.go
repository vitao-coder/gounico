package internal

import (
	"gounico/config"
	"gounico/pkg/database/dynamodb"
	"gounico/pkg/database/dynamodb/client"
	"gounico/pkg/logging"
	"gounico/pkg/logging/zap"

	"go.uber.org/fx"
)

var PackagesModule = fx.Provide(
	NewLogger,
	NewDynamoClient,
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
