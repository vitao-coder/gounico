package builders

import (
	"gounico/feiralivre/domains"
)

type FeiraLivreBuilder struct {
	feiraLivre       *domains.Feira
	regioesGenericas []*domains.RegiaoGenerica
}

func NewFeiraLivreBuilder() *FeiraLivreBuilder {
	return &FeiraLivreBuilder{
		feiraLivre: &domains.Feira{},
	}
}

func (flb *FeiraLivreBuilder) WithFeira(id int, nome string, registro string, setCens string, areaP string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.Id = id
	feira.Nome = nome
	feira.Registro = registro
	feira.SetCens = setCens
	feira.AreaP = areaP

	return flb
}

func (flb *FeiraLivreBuilder) WithLocalizacao(latitude float64, longitude float64, logradouro string, numero string, bairro string, referencia string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.Localizacao = domains.NewLocalizacao(latitude,
		longitude,
		logradouro,
		numero,
		bairro,
		referencia)
	return flb
}

func (flb *FeiraLivreBuilder) WithDistrito(id int, descricao string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.Distrito = domains.Distrito{
		Id:        id,
		Descricao: descricao,
	}

	feira.IdDistrito = feira.Distrito.Id
	return flb
}

func (flb *FeiraLivreBuilder) WithSubPrefeitura(id int, descricao string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.SubPrefeitura = domains.SubPrefeitura{
		Id:            id,
		SubPrefeitura: descricao,
		Regioes:       []domains.Regiao{},
	}
	feira.IdSubPref = feira.SubPrefeitura.Id
	return flb
}

func (flb *FeiraLivreBuilder) WithRegioes(descricaoRegiao5 string, descricaoRegiao8 string) *FeiraLivreBuilder {
	feira := flb.feiraLivre

	regiao5 := domains.NewRegiao(domains.CodigoRegiao(5))
	regiaoGenerica5 := domains.NewRegiaoGenerica(descricaoRegiao5)
	regiao5.IdRegiaoGenerica = regiaoGenerica5.Id

	regiao8 := domains.NewRegiao(domains.CodigoRegiao(8))
	regiaoGenerica8 := domains.NewRegiaoGenerica(descricaoRegiao8)
	regiao8.IdRegiaoGenerica = regiaoGenerica8.Id

	feira.SubPrefeitura.Regioes = append(feira.SubPrefeitura.Regioes, *regiao5, *regiao8)
	flb.regioesGenericas = append(flb.regioesGenericas, regiaoGenerica5, regiaoGenerica8)

	return flb
}

func (flb *FeiraLivreBuilder) BuildRegiaoGenerica() []*domains.RegiaoGenerica {
	return flb.regioesGenericas
}

func (flb *FeiraLivreBuilder) Build() *domains.Feira {
	return flb.feiraLivre
}
