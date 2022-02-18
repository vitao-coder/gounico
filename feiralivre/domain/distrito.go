package domain

type Distrito struct {
	Id        int    `gorm:"primary_key;column:ID"`
	Descricao string `gorm:"not null;"`
}
