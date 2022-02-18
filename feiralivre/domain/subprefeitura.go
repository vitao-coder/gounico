package domain

type SubPrefeitura struct {
	Id            int    `gorm:"primary_key;column:ID"`
	SubPrefeitura string `gorm:"column:SUBPREF"`
	IdRegiao      int    `gorm:"column:IDREGIAO"`
	Regiao        Regiao `gorm:"foreignkey:IdRegiaoGenerica;references:IdRegiao"`
}
