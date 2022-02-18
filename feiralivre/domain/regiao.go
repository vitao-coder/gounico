package domain

import (
	"crypto/md5"
	"fmt"

	"gorm.io/gorm"
)

type CodigoRegiao int

type Regiao struct {
	gorm.Model
	IdRegiaoGenerica string           `gorm:"type:char(32);index:,unique"`
	Regioes          []RegiaoGenerica `gorm:"many2many:regiao_regioes;foreignKey:IdRegiaoGenerica;References:IdRegiao;"`
	Codigo           int
}

type RegiaoGenerica struct {
	IdRegiao  string `gorm:"type:char(32);primary_key; index:,unique"`
	Descricao string `gorm:"type:varchar(255);not null"`
}

func NewRegiao(descricao string, codigoRegiao CodigoRegiao) *Regiao {

	regiao := &Regiao{
		Regioes: []RegiaoGenerica{
			{
				Descricao: descricao,
			}},
		Codigo: int(codigoRegiao),
	}

	regiao.Regioes[0].IdRegiao = regiao.Regioes[0].uniqueID()

	return regiao
}

func (r RegiaoGenerica) HashCode() string {
	return fmt.Sprintf("%s", r.Descricao)
}

func (r RegiaoGenerica) uniqueID() string {
	data := []byte(r.HashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
