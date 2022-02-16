package service

import (
	entityDomain "gounico/feiralivre/domain"
	"gounico/loaddata/domain"

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

	return nil, nil
}
