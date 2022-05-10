package internal

import (
	"gounico/feiralivre"
	feiraliv "gounico/feiralivre/services"
	"gounico/internal/repository"

	"go.uber.org/fx"
)

var ServicesModule = fx.Provide(
	NewProcessCSVService,
	NewFeiraLivreService,
)

func NewProcessCSVService(repository repository.Repository) feiralivre.ProcessCSV {
	loadDataService := feiraliv.NewProcessCSV(repository)
	return loadDataService
}

func NewFeiraLivreService(repository repository.Repository) feiralivre.FeiraLivre {
	feiraLivreService := feiraliv.NewFeiraLivreService(repository)
	return feiraLivreService
}
