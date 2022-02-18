package service

import (
	entityDomain "gounico/feiralivre/domain"
	"gounico/feiralivre/domain/builder"
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

func Test_feiraLivre_wrapDomainToEntities(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name            string
		feirasLivresCSV []*domain.FeirasLivresCSV
		want            []*entityDomain.Feira
		wantErr         bool
	}{
		{
			name: "If pass a array with correct information, return entities mounted to insert on database",
			feirasLivresCSV: []*domain.FeirasLivresCSV{
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
			want:    expectedEntities(),
			wantErr: false,
		},
		{
			name: "If pass a invalid data return a error",
			feirasLivresCSV: []*domain.FeirasLivresCSV{
				{
					Id:         "1",
					Longitude:  "-AA",
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
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl := &feiraLivre{}
			got, err := fl.wrapDomainToEntities(tt.feirasLivresCSV)
			if (err != nil) != tt.wantErr {
				t.Errorf("wrapDomainToEntities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapDomainToEntities() got = %v, want %v", got, tt.want)
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

func expectedEntities() []*entityDomain.Feira {

	var feirasEntities []*entityDomain.Feira

	feiraCSV := &domain.FeirasLivresCSV{
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
	}

	builderFeira := builder.NewFeiraLivreBuilder()
	builderFeira.
		WithFeira(1, feiraCSV.NomeFeira, feiraCSV.Registro, feiraCSV.SetCens, feiraCSV.AreaP).
		WithDistrito(87, feiraCSV.Distrito).
		WithLocalizacao(-23558733, -46550164, feiraCSV.Logradouro, feiraCSV.Numero, feiraCSV.Bairro, feiraCSV.Referencia).
		WithSubPrefeitura(26, feiraCSV.SubPrefe)

	builderFeira.AddRegiao(feiraCSV.Regiao5, 5)
	builderFeira.AddRegiao(feiraCSV.Regiao8, 8)

	feirasEntities = append(feirasEntities, builderFeira.Build())

	return feirasEntities
}
