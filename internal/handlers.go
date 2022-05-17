package internal

import (
	"gounico/feiralivre"
	"gounico/feiralivre/handlers/DELETE/excluirfeira"
	"gounico/feiralivre/handlers/GET/buscarfeira"
	"gounico/feiralivre/handlers/POST/novafeira"
	"gounico/feiralivre/handlers/POST/processcsv"
	"gounico/feiralivre/handlers/PUT/alterarfeira"
	"gounico/pkg/logging"
	"gounico/pkg/messaging/pulsar"
	"gounico/pkg/telemetry"

	"go.uber.org/fx"
)

var HandlersModule = fx.Provide(
	NewProcessCSVHandler,
	NewBuscarFeiraHandler,
	NewExcluirFeiraHandler,
	NewNovaFeiraHandler,
	NewAlterarFeiraHandler,
	NewNovaFeiraPublisher,
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

func NewNovaFeiraHandler(service feiralivre.FeiraLivre, telemetry telemetry.OpenTelemetry) HandlerOutput {
	handlerEndpoint := novafeira.NewNovaFeiraHandler(service, telemetry)
	return HandlerOutput{
		Endpoint: handlerEndpoint,
	}
}

func NewNovaFeiraPublisher(pulsar pulsar.PulsarClient, logger logging.Logger, telemetry telemetry.OpenTelemetry) HandlerOutput {
	handlerEndpoint := novafeira.NovaFeiraPublisher(pulsar, logger, telemetry)
	return HandlerOutput{
		Endpoint: handlerEndpoint,
	}
}

func NewAlterarFeiraHandler(service feiralivre.FeiraLivre) HandlerOutput {
	handlerEndpoint := alterarfeira.NewAlteraFeiraHandler(service)
	return HandlerOutput{
		Endpoint: handlerEndpoint,
	}
}
