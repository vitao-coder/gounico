package builder

import (
	"gounico/feiralivre/domain"
)

type FeiraLivreBuilder struct {
	feiraLivre       *domain.Feira
	regioesGenericas []*domain.RegiaoGenerica
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
	feira.Localizacao = domain.NewLocalizacao(latitude,
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

	feira.IdDistrito = feira.Distrito.Id
	return flb
}

func (flb *FeiraLivreBuilder) WithSubPrefeitura(id int, descricao string) *FeiraLivreBuilder {
	feira := flb.feiraLivre
	feira.SubPrefeitura = domain.SubPrefeitura{
		Id:            id,
		SubPrefeitura: descricao,
		Regioes:       []domain.Regiao{},
	}
	feira.IdSubPref = feira.SubPrefeitura.Id
	return flb
}

func (flb *FeiraLivreBuilder) WithRegioes(descricaoRegiao5 string, descricaoRegiao8 string) *FeiraLivreBuilder {
	feira := flb.feiraLivre

	regiao5 := domain.NewRegiao(domain.CodigoRegiao(5))
	regiaoGenerica5 := domain.NewRegiaoGenerica(descricaoRegiao5)
	regiao5.IdRegiaoGenerica = regiaoGenerica5.IdRegiao

	regiao8 := domain.NewRegiao(domain.CodigoRegiao(8))
	regiaoGenerica8 := domain.NewRegiaoGenerica(descricaoRegiao8)
	regiao8.IdRegiaoGenerica = regiaoGenerica8.IdRegiao

	feira.SubPrefeitura.Regioes = append(feira.SubPrefeitura.Regioes, *regiao5, *regiao8)
	flb.regioesGenericas = append(flb.regioesGenericas, regiaoGenerica5, regiaoGenerica8)

	return flb
}

func (flb *FeiraLivreBuilder) BuildRegiaoGenerica() []*domain.RegiaoGenerica {
	return flb.regioesGenericas
}

func (flb *FeiraLivreBuilder) Build() *domain.Feira {
	return flb.feiraLivre
}
