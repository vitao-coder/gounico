package domain

type SubPrefeitura struct {
	Id               int      `gorm:"primary_key;column:ID"`
	SubPrefeitura    string   `gorm:"column:SUBPREF"`
	IdRegiaoGenerica string   `gorm:"column:IDREGIAO"`
	Regioes          []Regiao `gorm:"foreignkey:IdRegiao;references:IdRegiaoGenerica"`
}
