package service

import (
	"context"
	"gounico/feiralivre/domain"
	"gounico/pkg/errors"
	"gounico/repository"
)

type feiraLivre struct {
	repository repository.Repository
}

func NewFeiraLivreService(repository repository.Repository) *feiraLivre {
	return &feiraLivre{
		repository: repository,
	}
}

func (f *feiraLivre) NovaFeira(feira domain.Feira) {

}

func (f *feiraLivre) ExcluirFeira(feiraID uint) {

}

func (f *feiraLivre) AlterarFeira(feiraID uint) {

}

func (f *feiraLivre) BuscarFeiraPorBairro(ctx context.Context, bairro string) ([]*domain.Feira, *errors.ServiceError) {
	var feirasPorBairro []*domain.Feira

	db := f.repository.DB().WithContext(ctx)

	err := db.Joins("JOIN localizacaos on localizacaos.uid=feiras.uid_localizacao").
		Where("localizacaos.bairro like ?", "%"+bairro+"%").Preload("Localizacao").Find(&feirasPorBairro)

	if err.Error != nil {
		return nil, errors.InternalServerError("erro ao consultar por bairro", err.Error)
	}

	if feirasPorBairro == nil {
		return nil, errors.NotFoundError()
	}

	return feirasPorBairro, nil
}

func (f *feiraLivre) BuscarFeiraPorDistrito(distrito string) {

}
