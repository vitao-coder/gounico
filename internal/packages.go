package internal

import (
	"gounico/pkg/logging"
	"gounico/pkg/logging/zap"

	"go.uber.org/fx"
)

var PackagesModule = fx.Provide(
	NewLogger,
)

func NewLogger() (logging.Logger, error) {
	logger, err := zap.NewZapLogger()
	return logger, err
}
