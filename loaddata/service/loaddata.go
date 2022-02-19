package service

import (
	"context"
	entityDomain "gounico/feiralivre/domain"
	"gounico/feiralivre/domain/builder"
	"gounico/loaddata/domain"
	"gounico/pkg/errors"
	"gounico/repository"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

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

	feiraLivreEntities, regioes, err := fl.wrapDomainToEntities(csvDomain)
	if err != nil {
		return errors.InternalServerError("Error translate CSV into domains.", err)
	}

	regioesDistincted := fl.distinctReusableData(regioes)
	err = fl.saveReusableData(ctx, regioesDistincted)
	if err != nil {
		return errors.InternalServerError("Error save reusable data.", err)
	}

	err = fl.saveDataToDatabase(ctx, feiraLivreEntities)
	if err != nil {
		return errors.InternalServerError("Error save reusable data.", err)
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

func (fl *loadData) wrapDomainToEntities(feirasLivresCSV []*domain.FeirasLivresCSV) ([]*entityDomain.Feira, []*entityDomain.RegiaoGenerica, error) {

	var feiraEntities []*entityDomain.Feira
	var regioes []*entityDomain.RegiaoGenerica

	for _, feiraCSV := range feirasLivresCSV {

		feiraID, distritoID, longitude, latitude, subPrefID, err := fl.convertStringsToBasicTypes(feiraCSV)
		if err != nil {
			return nil, nil, err
		}

		builderFeira := builder.NewFeiraLivreBuilder()
		builderFeira.
			WithFeira(feiraID, feiraCSV.NomeFeira, feiraCSV.Registro, feiraCSV.SetCens, feiraCSV.AreaP).
			WithDistrito(distritoID, feiraCSV.Distrito).
			WithLocalizacao(latitude, longitude, feiraCSV.Logradouro, feiraCSV.Numero, feiraCSV.Bairro, feiraCSV.Referencia).
			WithSubPrefeitura(subPrefID, feiraCSV.SubPrefe)

		builderFeira.WithRegioes(strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao5, " "), " "), strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao8, " "), " "))

		feiraEntity := builderFeira.Build()
		regioesGenericas := builderFeira.BuildRegiaoGenerica()
		regioes = append(regioes, regioesGenericas...)

		feiraEntities = append(feiraEntities, feiraEntity)

	}

	return feiraEntities, regioes, nil
}

func (fl *loadData) convertStringsToBasicTypes(feiraCSV *domain.FeirasLivresCSV) (feiraID int, distritoID int, longitude float64, latitude float64, subPrefID int, err error) {

	feiraID, err = strconv.Atoi(feiraCSV.Id)
	if err != nil {
		return
	}

	distritoID, err = strconv.Atoi(feiraCSV.CodDist)
	if err != nil {
		return
	}

	longitude, err = strconv.ParseFloat(feiraCSV.Longitude, 64)
	if err != nil {
		return
	}

	latitude, err = strconv.ParseFloat(feiraCSV.Latitude, 64)
	if err != nil {
		return
	}

	subPrefID, err = strconv.Atoi(feiraCSV.CodSubPref)
	if err != nil {
		return
	}

	return
}

func (fl *loadData) distinctReusableData(regioesGenericas []*entityDomain.RegiaoGenerica) []entityDomain.RegiaoGenerica {
	uniqueRegions := make(map[string]entityDomain.RegiaoGenerica)
	var regioesDistincted []entityDomain.RegiaoGenerica

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

func (fl *loadData) saveReusableData(ctx context.Context, regioes []entityDomain.RegiaoGenerica) error {

	if err := fl.repository.DB().WithContext(ctx).Model(&entityDomain.RegiaoGenerica{}).Create(&regioes); err.Error != nil {
		return err.Error
	}

	return nil
}

func (fl *loadData) saveDataToDatabase(ctx context.Context, feirasLivresDataToLoad []*entityDomain.Feira) error {

	if err := fl.repository.DB().WithContext(ctx).Model(&entityDomain.Feira{}).Create(feirasLivresDataToLoad); err.Error != nil {
		return err.Error
	}

	return nil
}
