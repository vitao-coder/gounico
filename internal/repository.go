package internal

import (
	"context"
	"gounico/internal/repository"
	"gounico/pkg/database"
	"gounico/pkg/logging"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	NewRepository,
)

func NewRepository(dynamoClient database.Database, logger logging.Logger) repository.Repository {
	ctx := context.Background()
	logger.Info(ctx, "Open connection with database...", nil)

	repo := repository.NewRepository(dynamoClient)

	logger.Info(ctx, "Connection successfully...", nil)
	return repo
}
