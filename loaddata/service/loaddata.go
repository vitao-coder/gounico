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

	feiraLivreEntities, err := fl.wrapDomainToEntities(csvDomain)
	if err != nil {
		return errors.InternalServerError("Error translate CSV into domains.", err)
	}

	regioes, localizacoes := fl.distinctReusableData(feiraLivreEntities)
	err = fl.saveReusableData(ctx, regioes, localizacoes)
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

func (fl *loadData) wrapDomainToEntities(feirasLivresCSV []*domain.FeirasLivresCSV) ([]*entityDomain.Feira, error) {

	var feiraEntities []*entityDomain.Feira

	for _, feiraCSV := range feirasLivresCSV {

		feiraID, distritoID, longitude, latitude, subPrefID, err := fl.convertStringsToBasicTypes(feiraCSV)
		if err != nil {
			return nil, err
		}

		builderFeira := builder.NewFeiraLivreBuilder()
		builderFeira.
			WithFeira(feiraID, feiraCSV.NomeFeira, feiraCSV.Registro, feiraCSV.SetCens, feiraCSV.AreaP).
			WithDistrito(distritoID, feiraCSV.Distrito).
			WithLocalizacao(latitude, longitude, feiraCSV.Logradouro, feiraCSV.Numero, feiraCSV.Bairro, feiraCSV.Referencia).
			WithSubPrefeitura(subPrefID, feiraCSV.SubPrefe)

		builderFeira.AddRegiao(strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao5, " "), " "), 5)
		builderFeira.AddRegiao(strings.TrimRight(strings.TrimLeft(feiraCSV.Regiao8, " "), " "), 8)

		feiraEntity := builderFeira.Build()

		feiraEntities = append(feiraEntities, feiraEntity)

	}

	return feiraEntities, nil
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

func (fl *loadData) distinctReusableData(feirasLivresDataToLoad []*entityDomain.Feira) ([]entityDomain.RegiaoGenerica, []entityDomain.Localizacao) {
	uniqueRegions := make(map[string]entityDomain.RegiaoGenerica)
	uniqueLocations := make(map[string]entityDomain.Localizacao)

	var regioesDistincted []entityDomain.RegiaoGenerica
	var localizacoesDistincted []entityDomain.Localizacao

	for _, feira := range feirasLivresDataToLoad {
		for _, regiao := range feira.SubPrefeitura.Regiao.Regioes {
			if _, ok := uniqueRegions[regiao.HashCode()]; !ok {
				uniqueRegions[regiao.Descricao] = regiao
			}
		}
		if _, ok := uniqueLocations[feira.Localizacao.HashCode()]; !ok {
			uniqueLocations[feira.Localizacao.UId] = feira.Localizacao
		}
	}

	for _, regionUnique := range uniqueRegions {
		regioesDistincted = append(regioesDistincted, regionUnique)
	}

	for _, localizacaoUnique := range uniqueLocations {
		localizacoesDistincted = append(localizacoesDistincted, localizacaoUnique)
	}
	return regioesDistincted, localizacoesDistincted
}

func (fl *loadData) saveReusableData(ctx context.Context, regioes []entityDomain.RegiaoGenerica, localizacoes []entityDomain.Localizacao) error {

	if err := fl.repository.DB().WithContext(ctx).Model(&entityDomain.RegiaoGenerica{}).Create(&regioes); err.Error != nil {
		return err.Error
	}

	if err := fl.repository.DB().WithContext(ctx).Model(&entityDomain.Localizacao{}).Create(&localizacoes); err.Error != nil {
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
