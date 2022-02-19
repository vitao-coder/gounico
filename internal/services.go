package internal

import (
	"gounico/feiralivre"
	feiralivreS "gounico/feiralivre/service"
	"gounico/loaddata"
	"gounico/loaddata/service"
	"gounico/repository"

	"go.uber.org/fx"
)

var ServicesModule = fx.Provide(
	NewLoadDataService,
	NewFeiraLivreService,
)

func NewLoadDataService(repository repository.Repository) loaddata.LoadData {
	loadDataService := service.NewLoadData(repository)
	return loadDataService
}

func NewFeiraLivreService(repository repository.Repository) feiralivre.FeiraLivre {
	feiraLivreService := feiralivreS.NewFeiraLivreService(repository)
	return feiraLivreService
}
