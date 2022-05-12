package internal

import (
	"gounico/feiralivre"
	"gounico/feiralivre/handlers/POST/novafeira"

	"go.uber.org/fx"
)

var ConsumersModule = fx.Provide(NewNovaFeiraConsumer)

type ConsumerOutput struct {
	fx.Out
	Consumer HTTPConsumer `group:"consumers"`
}

func NewNovaFeiraConsumer(service feiralivre.FeiraLivre) ConsumerOutput {
	consumerEndpoint := novafeira.NewNovaFeiraConsumer(service)
	return ConsumerOutput{
		Consumer: consumerEndpoint,
	}
}
