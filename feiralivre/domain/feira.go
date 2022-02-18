package domain

type Feira struct {
	Id            int           `gorm:"primary_key;column:ID"`
	Nome          string        `gorm:"not null;column:NOME"`
	Registro      string        `gorm:"not null;column:REGISTRO"`
	SetCens       string        `gorm:"not null;column:SETCENS"`
	AreaP         string        `gorm:"not null;column:AREAP"`
	IdLocalizacao string        `gorm:"not null;column:IDLOCA"`
	Localizacao   Localizacao   `gorm:"foreignKey:IdLocalizacao;references:UId"`
	IdDistrito    int           `gorm:"not null;column:IDDIST"`
	Distrito      Distrito      `gorm:"foreignKey:IdDistrito;references:Id"`
	IdSubPref     int           `gorm:"not null;column:IDSUBPREF"`
	SubPrefeitura SubPrefeitura `gorm:"foreignKey:IdSubPref;references:Id"`
}
