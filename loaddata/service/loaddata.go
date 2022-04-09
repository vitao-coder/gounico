package service

import (
	"context"
	domainFeira "gounico/feiralivre/domain"
	"gounico/feiralivre/domain/builder"
	"gounico/internal/repository"
	"gounico/loaddata/domain"
	"gounico/pkg/errors"
	"strings"

	"github.com/gocarina/gocsv"
)

const primaryType = "feira"
const secondaryType = "distrito"
const regiaoCutSet = " "

type loadData struct {
	repository repository.Repository
}

func NewLoadData(repository repository.Repository) *loadData {
	return &loadData{
		repository: repository,
	}
}

func (fl *loadData) ProcessCSVToDatabase(ctx context.Context, csvByteArray []byte) *errors.ServiceError {

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

func (fl *loadData) wrapCSVToDomain(csvByteArray []byte) ([]*domain.FeirasLivresCSV, error) {

	feirasLivresCSV := []*domain.FeirasLivresCSV{}

	err := gocsv.UnmarshalBytes(csvByteArray, &feirasLivresCSV)

	if err != nil {
		return nil, err
	}

	return feirasLivresCSV, nil
}

func (fl *loadData) wrapDomainToEntities(feirasLivresCSV []*domain.FeirasLivresCSV) ([]*domainFeira.Feira, error) {

	var feiraEntities []*domainFeira.Feira

	for _, feiraCSV := range feirasLivresCSV {

		feiraID, distritoID, longitude, latitude, subPrefID, err := feiraCSV.StringsToPrimitiveTypes()
		if err != nil {
			return nil, err
		}

		builderFeira := builder.NewFeiraLivreBuilder()
		builderFeira.
			WithFeira(feiraID, feiraCSV.NomeFeira, feiraCSV.Registro, feiraCSV.SetCens, feiraCSV.AreaP).
			WithDistrito(distritoID, feiraCSV.Distrito).
			WithLocalizacao(latitude, longitude, feiraCSV.Logradouro, feiraCSV.Numero, feiraCSV.Bairro, feiraCSV.Referencia).
			WithSubPrefeitura(subPrefID, feiraCSV.SubPrefe)

		builderFeira.WithRegioes(strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao5, regiaoCutSet), regiaoCutSet), strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao8, regiaoCutSet), regiaoCutSet))

		feiraEntity := builderFeira.Build()
		feiraEntity.Indexes(feiraCSV.Id, primaryType, feiraCSV.CodDist, secondaryType)
		feiraEntity.Data(feiraEntity)
		feiraEntities = append(feiraEntities, feiraEntity)
	}

	return feiraEntities, nil
}

func (fl *loadData) distinctReusableData(regioesGenericas []*domainFeira.RegiaoGenerica) []domainFeira.RegiaoGenerica {
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

func (fl *loadData) saveDataToDatabase(ctx context.Context, feirasLivresDataToLoad []*domainFeira.Feira) *errors.ServiceError {

	for _, feira := range feirasLivresDataToLoad {
		if err := fl.repository.Save(feira); err != nil {
			return err
		}
	}
	return nil
}
