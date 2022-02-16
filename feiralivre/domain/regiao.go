package domain

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type CodigoRegiao int

const (
	Regiao5 CodigoRegiao = 5
	Regiao8              = 8
)

type Regiao struct {
	Id             int
	RegiaoGenerica RegiaoGenerica
	Codigo         CodigoRegiao
}

type RegiaoGenerica struct {
	UId       string
	Descricao string
}

func NewRegiao(descricao string, codigoRegiao CodigoRegiao) *Regiao {

	regiao := &Regiao{
		RegiaoGenerica: RegiaoGenerica{
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
	uid := sha1.New()
	_, err := io.WriteString(uid, r.hashCode())
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s", uid.Sum(nil))
}
