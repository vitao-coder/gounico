package builder

import (
	"gounico/feiralivre/domain"
)

type FeiraLivreBuilder struct {
	feiraLivre *domain.Feira
}

func NewFeiraLivreBuilder() *FeiraLivreBuilder {
	return &FeiraLivreBuilder{
		feiraLivre: &domain.Feira{},
	}
}

func (flb *FeiraLivreBuilder) WithFeira(id int, nome string, registro string, setCens string, areaP string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.ID = uint(id)
	feira.Nome = nome
	feira.Registro = registro
	feira.SetCens = setCens
	feira.AreaP = areaP

	return flb
}

func (flb *FeiraLivreBuilder) WithLocalizacao(latitude float64, longitude float64, logradouro string, numero string, bairro string, referencia string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.Localizacao = *domain.NewLocalizacao(latitude,
		longitude,
		logradouro,
		numero,
		bairro,
		referencia)
	return flb
}

func (flb *FeiraLivreBuilder) WithDistrito(id int, descricao string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.Distrito = domain.Distrito{
		Id:        id,
		Descricao: descricao,
	}
	return flb
}

func (flb *FeiraLivreBuilder) WithSubPrefeitura(id int, descricao string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.SubPrefeitura = domain.SubPrefeitura{
		Id:            id,
		SubPrefeitura: descricao,
		Regioes:       []domain.Regiao{},
	}
	return flb
}

func (flb *FeiraLivreBuilder) AddRegiao(descricao string, codigoRegiao int) *FeiraLivreBuilder {
	feira := flb.feiraLivre

	regiao := domain.NewRegiao(descricao, domain.CodigoRegiao(codigoRegiao))

	feira.SubPrefeitura.Regioes = append(feira.SubPrefeitura.Regioes, *regiao)

	return flb
}

func (flb *FeiraLivreBuilder) Build() *domain.Feira {
	return flb.feiraLivre
}
