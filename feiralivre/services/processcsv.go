package services

import (
	"context"
	"gounico/constants"
	domainFeira "gounico/feiralivre/domains"
	"gounico/feiralivre/domains/builders"
	"gounico/internal/repository"
	"gounico/pkg/errors"
	"strings"

	"github.com/gocarina/gocsv"
)

type processCSV struct {
	repository repository.Repository
}

func NewProcessCSV(repository repository.Repository) *processCSV {
	return &processCSV{
		repository: repository,
	}
}

func (fl *processCSV) ProcessCSVToDatabase(ctx context.Context, csvByteArray []byte) *errors.ServiceError {

	csvDomain, err := fl.wrapCSVToDomain(csvByteArray)

	if err != nil {
		return errors.BadRequestError("Error wrap CSV. Please verify file columns and data types.")
	}

	feiraLivreEntities, err := fl.wrapDomainToEntities(csvDomain)
	if err != nil {
		return errors.InternalServerError("Error translate CSV into domains.", err)
	}

	errSave := fl.saveDataToDatabase(ctx, feiraLivreEntities)
	if errSave != nil {
		return errSave
	}

	return nil
}

func (fl *processCSV) wrapCSVToDomain(csvByteArray []byte) ([]*domainFeira.FeirasLivresCSV, error) {

	feirasLivresCSV := []*domainFeira.FeirasLivresCSV{}

	err := gocsv.UnmarshalBytes(csvByteArray, &feirasLivresCSV)

	if err != nil {
		return nil, err
	}

	return feirasLivresCSV, nil
}

func (fl *processCSV) wrapDomainToEntities(feirasLivresCSV []*domainFeira.FeirasLivresCSV) ([]*domainFeira.Feira, error) {

	var feiraEntities []*domainFeira.Feira

	for _, feiraCSV := range feirasLivresCSV {

		feiraID, distritoID, longitude, latitude, subPrefID, err := feiraCSV.StringsToPrimitiveTypes()
		if err != nil {
			return nil, err
		}

		builderFeira := builders.NewFeiraLivreBuilder()
		builderFeira.
			WithFeira(feiraID, feiraCSV.NomeFeira, feiraCSV.Registro, feiraCSV.SetCens, feiraCSV.AreaP).
			WithDistrito(distritoID, feiraCSV.Distrito).
			WithLocalizacao(latitude, longitude, feiraCSV.Logradouro, feiraCSV.Numero, feiraCSV.Bairro, feiraCSV.Referencia).
			WithSubPrefeitura(subPrefID, feiraCSV.SubPrefe)

		builderFeira.WithRegioes(strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao5, constants.RegiaoCutSet), constants.RegiaoCutSet), strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao8, constants.RegiaoCutSet), constants.RegiaoCutSet))

		feiraEntity := builderFeira.Build()
		feiraEntity.Indexes(feiraCSV.Id, constants.PrimaryType, feiraCSV.CodDist, constants.SecondaryType)
		feiraEntity.Data(feiraEntity)
		feiraEntities = append(feiraEntities, feiraEntity)
	}

	return feiraEntities, nil
}

func (fl *processCSV) distinctReusableData(regioesGenericas []*domainFeira.RegiaoGenerica) []domainFeira.RegiaoGenerica {
	uniqueRegions := make(map[string]domainFeira.RegiaoGenerica)
	var regioesDistincted []domainFeira.RegiaoGenerica

	for _, regiao := range regioesGenericas {
		if _, ok := uniqueRegions[regiao.HashCode()]; !ok {
			uniqueRegions[regiao.HashCode()] = *regiao
		}
	}
	for _, regionUnique := range uniqueRegions {
		regioesDistincted = append(regioesDistincted, regionUnique)
	}
	return regioesDistincted
}

func (fl *processCSV) saveDataToDatabase(ctx context.Context, feirasLivresDataToLoad []*domainFeira.Feira) *errors.ServiceError {
	var batchSaveArray []interface{}
	for _, feira := range feirasLivresDataToLoad {
		batchSaveArray = append(batchSaveArray, feira.DataDomain())
	}

	if err := fl.repository.BatchSave(batchSaveArray); err != nil {
		return err
	}
	return nil
}
