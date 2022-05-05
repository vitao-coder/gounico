package application

import "go.uber.org/fx"

var ConsumersModule = fx.Provide()

type ConsumerOutput struct {
	fx.Out
	Consumer HTTPEndpoint `group:"consumers"`
}
