package service

import (
	"gounico/loaddata/domain"

	"github.com/gocarina/gocsv"
)

type feiraLivre struct {
}

func NewFeiraLivre() *feiraLivre {
	return &feiraLivre{}
}

func (fl *feiraLivre) WrapCSVToDomain(csvByteArray []byte) ([]*domain.FeirasLivresCSV, error) {

	feirasLivresCSV := []*domain.FeirasLivresCSV{}

	err := gocsv.UnmarshalBytes(csvByteArray, &feirasLivresCSV)

	if err != nil {
		return nil, err
	}

	return feirasLivresCSV, nil
}
