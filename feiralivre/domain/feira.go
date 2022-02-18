package domain

import "gorm.io/gorm"

type Feira struct {
	gorm.Model
	Nome           string        `gorm:"not null"`
	Registro       string        `gorm:"not null"`
	SetCens        string        `gorm:"not null"`
	AreaP          string        `gorm:"not null"`
	UIdLocalizacao string        `gorm:"type:char(32); index:,unique"`
	Localizacao    Localizacao   `gorm:"foreignKey:UIdLocalizacao;References:UId;"`
	IdDistrito     int           `gorm:"not null"`
	Distrito       Distrito      `gorm:"foreignKey:IdDistrito;references:Id"`
	IdSubPref      int           `gorm:"not null;column:IDSUBPREF"`
	SubPrefeitura  SubPrefeitura `gorm:"foreignKey:IdSubPref;references:Id"`
}
