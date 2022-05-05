package domains

type SubPrefeitura struct {
	Id            int      `json:"id,omitempty"`
	SubPrefeitura string   `json:"sub_prefeitura,omitempty"`
	IdRegiao      int      `json:"id_regiao,omitempty"`
	Regioes       []Regiao `json:"regioes,omitempty"`
}
