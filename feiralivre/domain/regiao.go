package domain

import (
	"crypto/md5"
	"fmt"

	"gorm.io/gorm"
)

type CodigoRegiao int

const (
	Regiao5 CodigoRegiao = 5
	Regiao8 CodigoRegiao = 8
)

type Regiao struct {
	gorm.Model
	IdRegiaoGenerica string           `gorm:"column:REGIAO"`
	Regioes          []RegiaoGenerica `gorm:"many2many:regiao_regioes;association_foreignkey:UId;foreignkey:IdRegiaoGenerica"`
	Codigo           int
}

type RegiaoGenerica struct {
	UId       string `gorm:"primary_key;column:UID"`
	Descricao string
}

func NewRegiao(descricao string, codigoRegiao CodigoRegiao) *Regiao {

	regiao := &Regiao{
		Regioes: []RegiaoGenerica{
			{
				Descricao: descricao,
			}},
		Codigo: int(codigoRegiao),
	}

	return regiao
}

func (r Regiao) HashCode() string {
	return fmt.Sprintf("%s", r.Regioes[0].Descricao)
}

func (r Regiao) uniqueID() string {
	data := []byte(r.HashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
