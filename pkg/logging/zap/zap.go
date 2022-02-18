package zap

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logPath string) (*zapLogger, error) {

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		logPath,
	}

	logger, err := cfg.Build()

	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

func (zl zapLogger) Debug(ctx context.Context, msg string, object interface{}) {
	zl.logger.Debug(msg)
}

func (zl zapLogger) Info(ctx context.Context, msg string, object interface{}) {
	zl.logger.Info(msg)
}

func (zl zapLogger) Warn(ctx context.Context, msg string, object interface{}) {
	zl.logger.Warn(msg)
}

func (zl zapLogger) Error(ctx context.Context, msg string, obj interface{}, err error) {
	zl.logger.Error(msg, zap.Error(err))
}

func (zl zapLogger) Trace(ctx context.Context, begin time.Time, fc func() (command string, linesAffected int64), err error) {
	return
}
