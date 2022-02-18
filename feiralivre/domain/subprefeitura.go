package domain

type SubPrefeitura struct {
	Id            int      `gorm:"primary_key;column:ID"`
	SubPrefeitura string   `gorm:"primary_key;column:SUBPREF"`
	Regioes       []Regiao `gorm:"foreignKey:ID"`
}
