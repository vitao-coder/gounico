package services

import (
	"context"
	internalRepo "gounico/application/repository"
	"gounico/constants"
	"gounico/feiralivre/domains"
	"gounico/feiralivre/domains/builders"
	"gounico/infrastructure/errors"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type feiraLivre struct {
	repository internalRepo.Repository
}

func NewFeiraLivreService(repository internalRepo.Repository) *feiraLivre {
	return &feiraLivre{
		repository: repository,
	}
}

func (f *feiraLivre) SaveFeira(ctx context.Context, request *domains.FeiraRequest) *errors.ServiceError {
	feiraID, distritoID, longitude, latitude, subPrefID, err := request.StringsToPrimitiveTypes()
	if err != nil {
		return errors.BadRequestError("data request is not valid")
	}

	builderFeira := builders.NewFeiraLivreBuilder()
	builderFeira.
		WithFeira(feiraID, request.NomeFeira, request.Registro, request.SetCens, request.AreaP).
		WithDistrito(distritoID, request.Distrito).
		WithLocalizacao(latitude, longitude, request.Logradouro, request.Numero, request.Bairro, request.Referencia).
		WithSubPrefeitura(subPrefID, request.SubPrefe).
		WithRegioes(strings.TrimRight(strings.TrimLeft(request.Regiao5, constants.RegiaoCutSet), constants.RegiaoCutSet), strings.TrimRight(strings.TrimLeft(request.Regiao8, constants.RegiaoCutSet), constants.RegiaoCutSet))
	feiraEntity := builderFeira.Build()
	feiraEntity.Indexes(request.Id, constants.PrimaryType, request.CodDist, constants.SecondaryType)
	feiraEntity.Data(feiraEntity)

	if err := f.repository.Save(feiraEntity); err != nil {
		return err
	}
	return nil
}

func (f *feiraLivre) ExcluirFeira(ctx context.Context, feiraID string, distritoID string) *errors.ServiceError {

	if err := f.repository.Delete(feiraID, constants.PrimaryType, distritoID, constants.SecondaryType); err != nil {
		return err
	}

	return nil
}

func (f *feiraLivre) BuscarFeiraPorDistrito(ctx context.Context, distritoID string) ([]domains.Feira, *errors.ServiceError) {

	feiras, err := f.repository.GetBySecondaryID(distritoID, constants.SecondaryType)

	if err != nil {
		return nil, err
	}

	if feiras == nil {
		return nil, errors.NotFoundError()
	}

	feirasPorDistrito := []domains.Feira{}

	for _, domainEntity := range *feiras {
		var domainFeira domains.Feira
		if err := mapstructure.Decode(domainEntity.Data, &domainFeira); err != nil {
			return nil, errors.BadRequestError("Decoding data from BD error.")
		}
		feirasPorDistrito = append(feirasPorDistrito, domainFeira)
	}

	return feirasPorDistrito, nil
}