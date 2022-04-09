package service

import (
	"context"
	"gounico/feiralivre/domain"
	"gounico/feiralivre/domain/builder"
	internalRepo "gounico/internal/repository"
	"gounico/pkg/errors"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const primaryType = "feira"
const secondaryType = "distrito"
const regiaoCutSet = " "

type feiraLivre struct {
	repository internalRepo.Repository
}

func NewFeiraLivreService(repository internalRepo.Repository) *feiraLivre {
	return &feiraLivre{
		repository: repository,
	}
}

func (f *feiraLivre) NovaFeira(ctx context.Context, request *domain.FeiraRequest) *errors.ServiceError {
	feiraID, distritoID, longitude, latitude, subPrefID, err := request.StringsToPrimitiveTypes()
	if err != nil {
		return errors.BadRequestError("data request is not valid")
	}

	builderFeira := builder.NewFeiraLivreBuilder()
	builderFeira.
		WithFeira(feiraID, request.NomeFeira, request.Registro, request.SetCens, request.AreaP).
		WithDistrito(distritoID, request.Distrito).
		WithLocalizacao(latitude, longitude, request.Logradouro, request.Numero, request.Bairro, request.Referencia).
		WithSubPrefeitura(subPrefID, request.SubPrefe).
		WithRegioes(strings.TrimRight(strings.TrimLeft(request.Regiao5, regiaoCutSet), regiaoCutSet), strings.TrimRight(strings.TrimLeft(request.Regiao8, regiaoCutSet), regiaoCutSet))
	feiraEntity := builderFeira.Build()
	feiraEntity.Indexes(request.Id, primaryType, request.CodDist, secondaryType)
	feiraEntity.Data(feiraEntity)

	if err := f.repository.Save(feiraEntity); err != nil {
		return err
	}
	return nil
}

func (f *feiraLivre) ExcluirFeira(ctx context.Context, feiraID string) *errors.ServiceError {

	feira, err := f.repository.GetByPrimaryID(feiraID, primaryType)

	if err != nil {
		return err
	}
	if feira == nil {
		return errors.NotFoundError()
	}

	if err := f.repository.Delete(feiraID); err != nil {
		return err
	}

	return nil
}

func (f *feiraLivre) BuscarFeiraPorDistrito(ctx context.Context, distritoID string) ([]domain.Feira, *errors.ServiceError) {

	feiras, err := f.repository.GetBySecondaryID(distritoID, secondaryType)

	if err != nil {
		return nil, err
	}

	if feiras == nil {
		return nil, errors.NotFoundError()
	}

	feirasPorDistrito := []domain.Feira{}

	for _, domainEntity := range *feiras {
		var domainFeira domain.Feira
		if err := mapstructure.Decode(domainEntity.Data, &domainFeira); err != nil {
			return nil, errors.BadRequestError("Decoding data from BD error.")
		}
		feirasPorDistrito = append(feirasPorDistrito, domainFeira)
	}

	return feirasPorDistrito, nil
}
