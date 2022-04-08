package domain

import (
	"errors"
	"gounico/pkg/dynamodb/domain"
)

type Feira struct {
	*domain.DynamoDomain
	Id             int           `json:"id,omitempty"`
	Nome           string        `json:"nome,omitempty"`
	Registro       string        `json:"registro,omitempty"`
	SetCens        string        `json:"set_cens,omitempty"`
	AreaP          string        `json:"area_p,omitempty"`
	UIdLocalizacao string        `json:"u_id_localizacao,omitempty"`
	Localizacao    Localizacao   `json:"localizacao"`
	IdDistrito     int           `json:"id_distrito,omitempty"`
	Distrito       Distrito      `json:"distrito"`
	IdSubPref      int           `json:"id_sub_pref,omitempty"`
	SubPrefeitura  SubPrefeitura `json:"sub_prefeitura"`
}

func (f *Feira) IsDataValid() error {
	if f.ID == "" || f.PartitionID == "" || f.IDType == "" || f.PartitionIDType == "" {
		return errors.New("sensible data not informed")
	}
	return nil
}

func (f *Feira) DataDomain() interface{} {
	return f
}
