package zap

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (*zapLogger, error) {

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"../gounico/logs/gounico.log",
	}

	logger, err := cfg.Build()

	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

func (zl zapLogger) Debug(msg string) {
	zl.logger.Debug(msg)
}

func (zl zapLogger) Info(msg string) {
	zl.logger.Info(msg)
}

func (zl zapLogger) Warn(msg string) {
	zl.logger.Warn(msg)
}

func (zl zapLogger) Error(msg string, err error) {
	zl.logger.Error(msg, zap.Error(err))
}
