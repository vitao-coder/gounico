package internal

import (
	"gounico/loaddata"
	"gounico/loaddata/service"
	"gounico/repository"

	"go.uber.org/fx"
)

var ServicesModule = fx.Provide(
	NewLoadDataService,
)

func NewLoadDataService(repository repository.Repository) loaddata.LoadData {
	loadDataService := service.NewLoadData(repository)
	return loadDataService
}
