package service

import (
	"gounico/loaddata/domain"
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_feiraLivre_WrapCSVToDomain(t *testing.T) {
	tests := []struct {
		name         string
		csvByteArray []byte
		want         []*domain.FeirasLivresCSV
		wantErr      bool
	}{
		{
			name:         "If pass a valid CSV return a full converted domain FeirasLivresCSV",
			csvByteArray: readBytesFromCSVTestFile(),
			wantErr:      false,
			want: []*domain.FeirasLivresCSV{
				{
					Id:         "1",
					Longitude:  "-46550164",
					Latitude:   "-23558733",
					SetCens:    "355030885000091",
					AreaP:      "3550308005040",
					CodDist:    "87",
					Distrito:   "VILA FORMOSA",
					CodSubPref: "26",
					SubPrefe:   "ARICANDUVA-FORMOSA-CARRAO",
					Regiao5:    "Leste",
					Regiao8:    "Leste 1",
					NomeFeira:  "VILA FORMOSA",
					Registro:   "4041-0",
					Logradouro: "RUA MARAGOJIPE",
					Numero:     "S/N",
					Bairro:     "VL FORMOSA",
					Referencia: "TV RUA PRETORIA",
				},
			},
		},
		{
			name:         "If pass a invalid CSV return a error",
			csvByteArray: []byte{},
			wantErr:      true,
			want:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl := &feiraLivre{}
			got, err := fl.wrapCSVToDomain(tt.csvByteArray)
			if (err != nil) != tt.wantErr {
				t.Errorf("WrapCSVToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WrapCSVToDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func readBytesFromCSVTestFile() []byte {
	testBytes, err := ioutil.ReadFile("feiralivre_test.csv")
	if err != nil {
		panic(err)
	}
	return testBytes
}
