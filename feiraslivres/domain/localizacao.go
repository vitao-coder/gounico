package domain

type Localizacao struct {
	UId        string  `json:"id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Logradouro string  `json:"logradouro"`
	Numero     string  `json:"numero"`
	Bairro     string  `json:"bairro"`
	Referencia string  `json:"referencia"`
}
