package domain

import (
	"crypto/md5"
	"fmt"
)

type Localizacao struct {
	UId        string
	Latitude   float64
	Longitude  float64
	Logradouro string
	Numero     string
	Bairro     string
	Referencia string
}

func NewLocalizacao(latitude float64, longitude float64, logradouro string, numero string, bairro string, referencia string) *Localizacao {
	localizacao := &Localizacao{
		Latitude:   latitude,
		Longitude:  longitude,
		Logradouro: logradouro,
		Numero:     numero,
		Bairro:     bairro,
		Referencia: referencia,
	}
	localizacao.UId = localizacao.uniqueID()

	return localizacao
}

func (l Localizacao) hashCode() string {
	return fmt.Sprintf("%f/%f/%s/%s/%s", l.Latitude, l.Longitude, l.Logradouro, l.Numero, l.Bairro)
}

func (l Localizacao) uniqueID() string {
	data := []byte(l.hashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
