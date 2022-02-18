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
	IdRegiaoGenerica string
	RegiaoGenerica   RegiaoGenerica `gorm:"foreignkey:IdRegiaoGenerica;references:IdRegiaoGenerica"`
	Codigo           int
}

type RegiaoGenerica struct {
	gorm.Model
	IdRegiaoGenerica string
	Descricao        string
}

func NewRegiao(descricao string, codigoRegiao CodigoRegiao) *Regiao {

	regiao := &Regiao{
		RegiaoGenerica: RegiaoGenerica{
			Descricao: descricao,
		},
		Codigo: int(codigoRegiao),
	}
	regiao.RegiaoGenerica.IdRegiaoGenerica = regiao.uniqueID()
	regiao.IdRegiaoGenerica = regiao.RegiaoGenerica.IdRegiaoGenerica

	return regiao
}

func (r Regiao) HashCode() string {
	return fmt.Sprintf("%s", r.RegiaoGenerica.Descricao)
}

func (r Regiao) uniqueID() string {
	data := []byte(r.HashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
