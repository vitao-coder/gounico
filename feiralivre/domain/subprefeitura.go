package domain

type SubPrefeitura struct {
	Id            int      `gorm:"primary_key;"`
	SubPrefeitura string   `gorm:"not null"`
	IdRegiao      int      `gorm:"not null;index:,"`
	Regioes       []Regiao `gorm:"foreignKey:ID;References:IdRegiao;"`
}
