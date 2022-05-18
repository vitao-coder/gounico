package services

import (
	"context"
	"fmt"
	"gounico/feiralivre/domains"
	"gounico/feiralivre/domains/builders"
	"gounico/global"
	internalRepo "gounico/internal/repository"
	"gounico/pkg/errors"
	"gounico/pkg/telemetry/openTelemetry"
	"net/http"
	"strings"

	"github.com/google/uuid"

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
	ctx, traceSpan := openTelemetry.NewSpan(ctx, "Service - SaveFeira")
	defer traceSpan.End()
	feiraID, distritoID, longitude, latitude, subPrefID, err := request.StringsToPrimitiveTypes()
	if err != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
		openTelemetry.AddSpanError(traceSpan, err)
		return errors.BadRequestError("data request is not valid")
	}

	builderFeira := builders.NewFeiraLivreBuilder()
	builderFeira.
		WithFeira(feiraID, request.NomeFeira, request.Registro, request.SetCens, request.AreaP).
		WithDistrito(distritoID, request.Distrito).
		WithLocalizacao(latitude, longitude, request.Logradouro, request.Numero, request.Bairro, request.Referencia).
		WithSubPrefeitura(subPrefID, request.SubPrefe).
		WithRegioes(strings.TrimRight(strings.TrimLeft(request.Regiao5, global.RegiaoCutSet), global.RegiaoCutSet), strings.TrimRight(strings.TrimLeft(request.Regiao8, global.RegiaoCutSet), global.RegiaoCutSet))
	feiraEntity := builderFeira.Build()
	request.Id = uuid.New().String()
	feiraEntity.Indexes(request.Id, global.PrimaryType, request.CodDist, global.SecondaryType)
	feiraEntity.Data(feiraEntity)
	indexes := feiraEntity.GetIndexes()

	openTelemetry.AddSpanTags(traceSpan, buildFeiraTelemetryTags(indexes.PartitionKey, indexes.SortKey, indexes.UID))
	if err := f.repository.Save(ctx, feiraEntity); err != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
		openTelemetry.AddSpanError(traceSpan, err)
		return err
	}

	openTelemetry.SuccessSpan(traceSpan, fmt.Sprintf("StatusCode: %d", http.StatusCreated))
	return nil
}

func (f *feiraLivre) ExcluirFeira(ctx context.Context, feiraID string, distritoID string) *errors.ServiceError {

	if err := f.repository.Delete(feiraID, global.PrimaryType, distritoID, global.SecondaryType); err != nil {
		return err
	}

	return nil
}

func (f *feiraLivre) BuscarFeiraPorDistrito(ctx context.Context, distritoID string) ([]domains.Feira, *errors.ServiceError) {

	feiras, err := f.repository.GetBySecondaryID(distritoID, global.SecondaryType)

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

func buildFeiraTelemetryTags(partitionKey string, sortKey string, uid string) map[string]string {
	return map[string]string{
		"partitionKey": partitionKey,
		"sortKey":      sortKey,
		"UID":          uid,
	}
}
