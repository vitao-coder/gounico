package service

import (
	"context"
	"gounico/feiralivre/domain"
	"gounico/feiralivre/domain/builder"
	"gounico/pkg/errors"
	"gounico/repository"
	"strings"
)

type feiraLivre struct {
	repository repository.Repository
}

func NewFeiraLivreService(repository repository.Repository) *feiraLivre {
	return &feiraLivre{
		repository: repository,
	}
}

func (f *feiraLivre) NovaFeira(ctx context.Context, request *domain.FeiraRequest) *errors.ServiceError {
	feiraID, distritoID, longitude, latitude, subPrefID, err := request.StringsToPrimitiveTypes()
	if err != nil {
		return nil
	}

	builderFeira := builder.NewFeiraLivreBuilder()
	builderFeira.
		WithFeira(feiraID, request.NomeFeira, request.Registro, request.SetCens, request.AreaP).
		WithDistrito(distritoID, request.Distrito).
		WithLocalizacao(latitude, longitude, request.Logradouro, request.Numero, request.Bairro, request.Referencia).
		WithSubPrefeitura(subPrefID, request.SubPrefe)

	builderFeira.WithRegioes(strings.TrimRight(strings.TrimLeft(request.Regiao5, " "), " "), strings.TrimRight(strings.TrimLeft(request.Regiao8, " "), " "))
	feiraEntity := builderFeira.Build()

	localizacao := *&domain.Localizacao{}

	db := f.repository.DB()

	if locaErr := db.Model(&domain.Localizacao{}).First(&localizacao, "uid = ?", feiraEntity.Localizacao.UId); locaErr.Error != nil {
		if locaErr.Error.Error() != "record not found" {
			return errors.InternalServerError("Error save new Feira", locaErr.Error)
		}
	}
	dbFeira := db.WithContext(ctx).Model(&domain.Feira{})
	if localizacao.UId != "" {
		feiraEntity.UIdLocalizacao = localizacao.UId
		dbFeira = dbFeira.Omit("Localizacao")
	}

	if err := dbFeira.Create(&feiraEntity); err.Error != nil {
		return errors.InternalServerError("Error save new Feira", err.Error)
	}
	return nil
}

func (f *feiraLivre) ExcluirFeira(ctx context.Context, feiraID uint) *errors.ServiceError {
	db := f.repository.DB().WithContext(ctx)

	var exists bool
	err := db.Model(&domain.Feira{}).
		Select("count(*) > 0").
		Where("id = ?", feiraID).
		Find(&exists)

	if !exists {
		return errors.NotFoundError()
	}

	err = db.Delete(&domain.Feira{}, feiraID)

	if err.Error != nil {
		return errors.InternalServerError("Error on delete.", err.Error)
	}
	return nil
}

func (f *feiraLivre) AlterarFeira(feiraID uint, request domain.FeiraRequest) {

}

func (f *feiraLivre) BuscarFeiraPorBairro(ctx context.Context, bairro string) ([]*domain.Feira, *errors.ServiceError) {
	var feirasPorBairro []*domain.Feira

	db := f.repository.DB().WithContext(ctx)

	err := db.Joins("JOIN localizacaos on localizacaos.uid=feiras.uid_localizacao").
		Where("localizacaos.bairro like ?", "%"+bairro+"%").
		Preload("Localizacao").
		Preload("Distrito").
		Preload("SubPrefeitura").
		Find(&feirasPorBairro)

	if err.Error != nil {
		return nil, errors.InternalServerError("erro ao consultar por bairro", err.Error)
	}

	if feirasPorBairro == nil || len(feirasPorBairro) < 1 {
		return nil, errors.NotFoundError()
	}

	return feirasPorBairro, nil
}

func (f *feiraLivre) BuscarFeiraPorDistrito(ctx context.Context, distrito string) ([]*domain.Feira, *errors.ServiceError) {
	var feirasPorDistrito []*domain.Feira

	db := f.repository.DB().WithContext(ctx)

	err := db.Joins("JOIN distritos on distritos.ID=feiras.id_distrito").
		Where("distritos.descricao like ?", "%"+distrito+"%").
		Preload("Localizacao").
		Preload("Distrito").
		Preload("SubPrefeitura").
		Find(&feirasPorDistrito)

	if err.Error != nil {
		return nil, errors.InternalServerError("erro ao consultar por distrito", err.Error)
	}

	if feirasPorDistrito == nil || len(feirasPorDistrito) < 1 {
		return nil, errors.NotFoundError()
	}

	return feirasPorDistrito, nil
}
