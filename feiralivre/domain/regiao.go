package domain

import (
	"crypto/md5"
	"fmt"

	"gorm.io/gorm"
)

type CodigoRegiao int

type Regiao struct {
	gorm.Model
	IdRegiaoGenerica string `gorm:"type:char(32);index:,unique"`
	Codigo           int
}

type RegiaoGenerica struct {
	IdRegiao  string   `gorm:"type:char(32);primary_key; index:,unique"`
	Descricao string   `gorm:"type:varchar(255);not null"`
	Regioes   []Regiao `gorm:"many2many:regiao_regioes;foreignKey:IdRegiao;References:IdRegiaoGenerica;"`
}

func NewRegiaoGenerica(descricao string) *RegiaoGenerica {
	regiao := &RegiaoGenerica{
		Descricao: descricao,
	}

	regiao.IdRegiao = regiao.uniqueID()
	return regiao
}

func NewRegiao(codigoRegiao CodigoRegiao) *Regiao {

	regiao := &Regiao{
		Codigo: int(codigoRegiao),
	}

	return regiao
}

func (r RegiaoGenerica) HashCode() string {
	return fmt.Sprintf("%s", r.Descricao)
}

func (r RegiaoGenerica) uniqueID() string {
	data := []byte(r.HashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
