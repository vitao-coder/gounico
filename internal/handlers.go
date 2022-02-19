package internal

import (
	"gounico/feiralivre"
	"gounico/feiralivre/handler/GET/buscarfeira"
	"gounico/loaddata"
	"gounico/loaddata/handler/POST/processcsv"

	"go.uber.org/fx"
)

var HandlersModule = fx.Provide(
	NewProcessCSVHandler,
	NewBuscarFeiraHandler,
)

type HandlerOutput struct {
	fx.Out
	Endpoint HTTPEndpoint `group:"endpoints"`
}

func NewProcessCSVHandler(service loaddata.LoadData) HandlerOutput {
	handlerEndpoint := processcsv.NewProcessCSVHandler(service)
	return HandlerOutput{
		Endpoint: handlerEndpoint,
	}
}

func NewBuscarFeiraHandler(service feiralivre.FeiraLivre) HandlerOutput {
	handlerEndpoint := buscarfeira.NewBuscarFeiraHandler(service)
	return HandlerOutput{
		Endpoint: handlerEndpoint,
	}
}
