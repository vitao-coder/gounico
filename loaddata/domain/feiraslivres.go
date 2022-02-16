package domain

type FeirasLivresCSV struct {
	Id         string `csv:"ID"`
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
