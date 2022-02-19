package domain

import (
	"crypto/md5"
	"fmt"
)

type Localizacao struct {
	UId        string  `gorm:"type:char(32);primary_key; index:,unique"`
	Latitude   float64 `gorm:"not null"`
	Longitude  float64 `gorm:"not null"`
	Logradouro string  `gorm:"not null"`
	Numero     string  `gorm:"not null"`
	Bairro     string  `gorm:"not null"`
	Referencia string  `gorm:"not null"`
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
