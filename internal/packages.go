package internal

import (
	"gounico/config"
	"gounico/pkg/logging"
	"gounico/pkg/logging/zap"

	"go.uber.org/fx"
)

var PackagesModule = fx.Provide(
	NewLogger,
)

func NewLogger(config config.Configuration) (logging.Logger, error) {
	logger, err := zap.NewZapLogger(config.Server.LogPath)
	return logger, err
}
