package domain

import "gorm.io/gorm"

type Feira struct {
	gorm.Model
	Nome           string        `gorm:"not null;column:NOME"`
	Registro       string        `gorm:"not null;column:REGISTRO"`
	SetCens        string        `gorm:"not null;column:SETCENS"`
	AreaP          string        `gorm:"not null;column:AREAP"`
	UIdLocalizacao string        `gorm:"not null;size:255;column:UIdLocalizacao"`
	Localizacao    Localizacao   `gorm:"foreignKey:UIdLocalizacao;references:UId"`
	IdDistrito     int           `gorm:"not null;column:IDDIST"`
	Distrito       Distrito      `gorm:"foreignKey:IdDistrito;references:Id"`
	IdSubPref      int           `gorm:"not null;column:IDSUBPREF"`
	SubPrefeitura  SubPrefeitura `gorm:"foreignKey:IdSubPref;references:Id"`
}
