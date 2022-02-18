package service

import (
	entityDomain "gounico/feiralivre/domain"
	"gounico/feiralivre/domain/builder"
	"gounico/loaddata/domain"
	"strconv"

	"github.com/gocarina/gocsv"
)

type feiraLivre struct {
}

func NewFeiraLivre() *feiraLivre {
	return &feiraLivre{}
}

func (fl *feiraLivre) wrapCSVToDomain(csvByteArray []byte) ([]*domain.FeirasLivresCSV, error) {

	feirasLivresCSV := []*domain.FeirasLivresCSV{}

	err := gocsv.UnmarshalBytes(csvByteArray, &feirasLivresCSV)

	if err != nil {
		return nil, err
	}

	return feirasLivresCSV, nil
}

func (fl *feiraLivre) wrapDomainToEntities(feirasLivresCSV []*domain.FeirasLivresCSV) ([]*entityDomain.Feira, error) {

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

		builderFeira.AddRegiao(feiraCSV.Regiao5, 5)
		builderFeira.AddRegiao(feiraCSV.Regiao8, 8)

		feiraEntity := builderFeira.Build()

		feiraEntities = append(feiraEntities, feiraEntity)

	}

	return feiraEntities, nil
}

func (fl *feiraLivre) convertStringsToBasicTypes(feiraCSV *domain.FeirasLivresCSV) (feiraID int, distritoID int, longitude float64, latitude float64, subPrefID int, err error) {

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

func (fl *feiraLivre) processDataToDatabase(feirasLivresDataToLoad []*entityDomain.Feira) (bool, error) {

	return true, nil
}

func (fl *feiraLivre) distinctReusableData(feirasLivresDataToLoad []*entityDomain.Feira) ([]entityDomain.RegiaoGenerica, []entityDomain.Localizacao) {
	return nil, nil
}

func (fl *feiraLivre) saveReusableData(feirasLivresDataToLoad []*entityDomain.Feira) ([]entityDomain.RegiaoGenerica, []entityDomain.Localizacao) {
	return nil, nil
}

func (fl *feiraLivre) saveDataToDatabase(feirasLivresDataToLoad []*entityDomain.Feira) (bool, error) {

	return true, nil
}
