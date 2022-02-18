package internal

import (
	"context"
	"gounico/config"
	"gounico/pkg/logging"
	"gounico/repository"
	"gounico/repository/database/gorm"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	NewRepository,
)

func NewRepository(config config.Configuration, logger logging.Logger) (repository.Repository, error) {
	ctx := context.Background()

	logger.Info(ctx, "Open connection with database...", nil)
	repo, err := gorm.NewDatabase(config.Database)
	if err != nil {
		logger.Error(ctx, "Error open connection with database.", nil, err)
		return nil, err
	}
	logger.Info(context.Background(), "Database connection successful.", nil)

	logger.Info(ctx, "Executing migrations for database...", nil)
	err = repo.AutoMigrateDatabase()
	if err != nil {
		logger.Error(ctx, "Error execute migrations for database.", nil, err)
		return nil, err
	}

	logger.Info(ctx, "Migrations executed successfully.", nil)

	return repo, nil
}
