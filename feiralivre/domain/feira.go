package domain

type Feira struct {
	Id            int
	Nome          string
	Registro      string
	SetCens       string
	AreaP         string
	Localizacao   *Localizacao
	Distrito      *Distrito
	SubPrefeitura *SubPrefeitura
}
