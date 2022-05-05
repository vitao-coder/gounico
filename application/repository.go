package application

import (
	"context"
	"gounico/application/repository"
	"gounico/infrastructure/database/dynamodb"
	"gounico/infrastructure/logging"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	NewRepository,
)

func NewRepository(dynamoClient dynamodb.DynamoClient, logger logging.Logger) repository.Repository {
	ctx := context.Background()
	logger.Info(ctx, "Open connection with database...", nil)

	repo := repository.NewRepository(dynamoClient)

	logger.Info(ctx, "Connection successfully...", nil)
	return repo
}
