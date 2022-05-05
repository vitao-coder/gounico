package domains

import (
	"strconv"
)

type FeiraRequest struct {
	Id         string `json:"id,omitempty" `
	Longitude  string `json:"longitude" validate:"required"`
	Latitude   string `json:"latitude" validate:"required"`
	SetCens    string `json:"set_cens" validate:"required"`
	AreaP      string `json:"area_p" validate:"required"`
	CodDist    string `json:"cod_dist" validate:"required"`
	Distrito   string `json:"distrito" validate:"required"`
	CodSubPref string `json:"cod_sub_pref" validate:"required"`
	SubPrefe   string `json:"sub_prefe" validate:"required"`
	Regiao5    string `json:"regiao_5" validate:"required"`
	Regiao8    string `json:"regiao_8" validate:"required"`
	NomeFeira  string `json:"nome_feira" validate:"required"`
	Registro   string `json:"registro" validate:"required"`
	Logradouro string `json:"logradouro" validate:"required"`
	Numero     string `json:"numero" validate:"required"`
	Bairro     string `json:"bairro" validate:"required"`
	Referencia string `json:"referencia" validate:"required"`
}

func (fr *FeiraRequest) StringsToPrimitiveTypes() (feiraID int, distritoID int, longitude float64, latitude float64, subPrefID int, err error) {

	feiraID, err = strconv.Atoi(fr.Id)
	if err != nil {
		return
	}

	distritoID, err = strconv.Atoi(fr.CodDist)
	if err != nil {
		return
	}

	longitude, err = strconv.ParseFloat(fr.Longitude, 64)
	if err != nil {
		return
	}

	latitude, err = strconv.ParseFloat(fr.Latitude, 64)
	if err != nil {
		return
	}

	subPrefID, err = strconv.Atoi(fr.CodSubPref)
	if err != nil {
		return
	}

	return
}
