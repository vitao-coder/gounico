package domain

import "strconv"

type FeirasLivresCSV struct {
	Id         string `csv:"UID"`
	Longitude  string `csv:"LONG"`
	Latitude   string `csv:"LAT"`
	SetCens    string `csv:"SETCENS"`
	AreaP      string `csv:"AREAP"`
	CodDist    string `csv:"CODDIST"`
	Distrito   string `csv:"DISTRITO"`
	CodSubPref string `csv:"CODSUBPREF"`
	SubPrefe   string `csv:"SUBPREFE"`
	Regiao5    string `csv:"REGIAO5"`
	Regiao8    string `csv:"REGIAO8"`
	NomeFeira  string `csv:"NOME_FEIRA"`
	Registro   string `csv:"REGISTRO"`
	Logradouro string `csv:"LOGRADOURO"`
	Numero     string `csv:"NUMERO"`
	Bairro     string `csv:"BAIRRO"`
	Referencia string `csv:"REFERENCIA"`
}

func (flcsv *FeirasLivresCSV) StringsToPrimitiveTypes() (feiraID int, distritoID int, longitude float64, latitude float64, subPrefID int, err error) {

	feiraID, err = strconv.Atoi(flcsv.Id)
	if err != nil {
		return
	}

	distritoID, err = strconv.Atoi(flcsv.CodDist)
	if err != nil {
		return
	}

	longitude, err = strconv.ParseFloat(flcsv.Longitude, 64)
	if err != nil {
		return
	}

	latitude, err = strconv.ParseFloat(flcsv.Latitude, 64)
	if err != nil {
		return
	}

	subPrefID, err = strconv.Atoi(flcsv.CodSubPref)
	if err != nil {
		return
	}

	return
}
