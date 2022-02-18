package domain

import (
	"crypto/md5"
	"fmt"
)

type Localizacao struct {
	UId        string  `gorm:"primary_key;column:UID"`
	Latitude   float64 `gorm:"not null;column:LAT"`
	Longitude  float64 `gorm:"not null;column:LONG"`
	Logradouro string  `gorm:"not null;column:LOGRA"`
	Numero     string  `gorm:"not null;column:NUMERO"`
	Bairro     string  `gorm:"not null;column:BAIRRO"`
	Referencia string  `gorm:"not null;column:REFERENCIA"`
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

func (l Localizacao) HashCode() string {
	return fmt.Sprintf("%f/%f/%s/%s/%s", l.Latitude, l.Longitude, l.Logradouro, l.Numero, l.Bairro)
}

func (l Localizacao) uniqueID() string {
	data := []byte(l.HashCode())
	return fmt.Sprintf("%x", md5.Sum(data))
}
