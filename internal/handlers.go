package internal

import (
	"gounico/feiralivre"
	"gounico/feiralivre/handler/DELETE/excluirfeira"
	"gounico/feiralivre/handler/GET/buscarfeira"
	"gounico/feiralivre/handler/POST/novafeira"
	"gounico/feiralivre/handler/POST/processcsv"

	"go.uber.org/fx"
)

var HandlersModule = fx.Provide(
	NewProcessCSVHandler,
	NewBuscarFeiraHandler,
	NewExcluirFeiraHandler,
	NewNovaFeiraHandler,
)

type HandlerOutput struct {
	fx.Out
	Endpoint HTTPEndpoint `group:"endpoints"`
}

func NewProcessCSVHandler(service feiralivre.ProcessCSV) HandlerOutput {
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

func NewExcluirFeiraHandler(service feiralivre.FeiraLivre) HandlerOutput {
	handlerEndpoint := excluirfeira.NewExcluirFeiraHandler(service)
	return HandlerOutput{
		Endpoint: handlerEndpoint,
	}
}

func NewNovaFeiraHandler(service feiralivre.FeiraLivre) HandlerOutput {
	handlerEndpoint := novafeira.NewNovaFeiraHandler(service)
	return HandlerOutput{
		Endpoint: handlerEndpoint,
	}
}
