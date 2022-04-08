package domain

import (
	"crypto/md5"
	"fmt"
)

type RegiaoGenerica struct {
	Id        string   `json:"id,omitempty"`
	Descricao string   `json:"descricao,omitempty"`
	Regioes   []Regiao `json:"regioes,omitempty"`
}

func NewRegiaoGenerica(descricao string) *RegiaoGenerica {
	regiao := &RegiaoGenerica{
		Descricao: descricao,
	}

	regiao.Id = regiao.uniqueID()
	return regiao
}

func (r RegiaoGenerica) HashCode() string {
	return fmt.Sprintf("%s", r.Descricao)
}

func (r RegiaoGenerica) uniqueID() string {
	data := []byte(r.HashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
