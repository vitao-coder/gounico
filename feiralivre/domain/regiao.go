package domain

type CodigoRegiao int

type Regiao struct {
	IdRegiaoGenerica string `json:"id_regiao_generica,omitempty"`
	Codigo           int    `json:"codigo,omitempty"`
}

func NewRegiao(codigoRegiao CodigoRegiao) *Regiao {

	regiao := &Regiao{
		Codigo: int(codigoRegiao),
	}

	return regiao
}
