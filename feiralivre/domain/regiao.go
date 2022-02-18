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
	RegiaoGenerica *RegiaoGenerica `gorm:"foreignKey:UID"`
	Codigo         CodigoRegiao    `gorm:"not null;column:CODIGOREGIAO"`
}

type RegiaoGenerica struct {
	UId       string `gorm:"primary_key;column:UID"`
	Descricao string `gorm:"not null;column:DESCRICAO"`
}

func NewRegiao(descricao string, codigoRegiao CodigoRegiao) *Regiao {

	regiao := &Regiao{
		RegiaoGenerica: &RegiaoGenerica{
			Descricao: descricao,
		},
		Codigo: codigoRegiao,
	}
	regiao.RegiaoGenerica.UId = regiao.uniqueID()

	return regiao
}

func (r Regiao) hashCode() string {
	return fmt.Sprintf("%s", r.RegiaoGenerica.Descricao)
}

func (r Regiao) uniqueID() string {
	data := []byte(r.hashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
