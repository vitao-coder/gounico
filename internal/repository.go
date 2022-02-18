package internal

import (
	"gounico/config"
	"gounico/repository"
	"gounico/repository/database/gorm"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	NewBaseRepository,
)

func NewBaseRepository(config config.Configuration) (repository.Repository, error) {
	repository, err := gorm.NewDatabase(config.Server.DBName)
	if err != nil {
		return nil, err
	}
	return repository, nil
}
