package internal

import (
	"gounico/loaddata"
	"gounico/loaddata/handler/POST/processcsv"

	"go.uber.org/fx"
)

var HandlersModule = fx.Provide(
	NewProcessCSVHandler,
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
