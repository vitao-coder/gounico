package internal

import (
	"gounico/feiralivre"
	feiraliv "gounico/feiralivre/service"
	"gounico/internal/repository"
	"gounico/loaddata"
	"gounico/loaddata/service"

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
	feiraLivreService := feiraliv.NewFeiraLivreService(repository)
	return feiraLivreService
}
