package domains

import (
	"crypto/md5"
	"fmt"
)

type Localizacao struct {
	Id         string  `json:"id,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Logradouro string  `json:"logradouro,omitempty"`
	Numero     string  `json:"numero,omitempty"`
	Bairro     string  `json:"bairro,omitempty"`
	Referencia string  `json:"referencia,omitempty"`
}

func NewLocalizacao(latitude float64, longitude float64, logradouro string, numero string, bairro string, referencia string) Localizacao {
	localizacao := Localizacao{
		Latitude:   latitude,
		Longitude:  longitude,
		Logradouro: logradouro,
		Numero:     numero,
		Bairro:     bairro,
		Referencia: referencia,
	}
	localizacao.Id = localizacao.uniqueID()

	return localizacao
}

func (l Localizacao) HashCode() string {
	return fmt.Sprintf("%f/%f/%s/%s/%s", l.Latitude, l.Longitude, l.Logradouro, l.Numero, l.Bairro)
}

func (l Localizacao) uniqueID() string {
	data := []byte(l.HashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
