package domain

type Feira struct {
	Id            int            `gorm:"primary_key;column:ID"`
	Nome          string         `gorm:"not null;column:NOME"`
	Registro      string         `gorm:"not null;column:REGISTRO"`
	SetCens       string         `gorm:"not null;column:SETCENS"`
	AreaP         string         `gorm:"not null;column:AREAP"`
	Localizacao   *Localizacao   `gorm:"foreignKey:UID"`
	Distrito      *Distrito      `gorm:"foreignKey:ID"`
	SubPrefeitura *SubPrefeitura `gorm:"foreignKey:ID"`
}
