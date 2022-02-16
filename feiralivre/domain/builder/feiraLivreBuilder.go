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

func (flb *FeiraLivreBuilder) WithFeira(id int, nome string, registro string, setCens string, areaP string) *domain.Feira {
	feira := flb.feiraLivre
	feira.Id = id
	feira.Nome = nome
	feira.Registro = registro
	feira.SetCens = setCens
	feira.AreaP = areaP

	return feira
}

func (flb *FeiraLivreBuilder) WithLocalizacao(latitude float64, longitude float64, logradouro string, numero string, bairro string, referencia string) *domain.Feira {
	feira := flb.feiraLivre
	feira.Localizacao = domain.NewLocalizacao(latitude,
		longitude,
		logradouro,
		numero,
		bairro,
		referencia)
	return feira
}

func (flb *FeiraLivreBuilder) WithDistrito(id int, descricao string) *domain.Feira {
	feira := flb.feiraLivre
	feira.Distrito = &domain.Distrito{
		Id:        id,
		Descricao: descricao,
	}
	return feira
}

func (flb *FeiraLivreBuilder) WithSubPrefeitura(id int, descricao string) *domain.Feira {
	feira := flb.feiraLivre
	feira.SubPrefeitura = &domain.SubPrefeitura{
		Id:            id,
		SubPrefeitura: descricao,
		Regioes:       []domain.Regiao{},
	}
	return feira
}

func (flb *FeiraLivreBuilder) AddRegiao(descricao string, codigoRegiao int) *domain.Feira {
	feira := flb.feiraLivre

	if feira.SubPrefeitura == nil {
		return nil
	}
	regiao := domain.NewRegiao(descricao, domain.CodigoRegiao(codigoRegiao))
	feira.SubPrefeitura.Regioes = append(feira.SubPrefeitura.Regioes, *regiao)

	return feira
}

func (flb *FeiraLivreBuilder) Build() *domain.Feira {
	return flb.feiraLivre
}
