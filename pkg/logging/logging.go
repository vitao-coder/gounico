package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Logger interface {
	GetInternalLogger() *zap.Logger
	Debug(ctx context.Context, msg string, obj interface{})
	Info(ctx context.Context, msg string, obj interface{})
	Warn(ctx context.Context, msg string, obj interface{})
	Error(ctx context.Context, msg string, obj interface{}, err error)
	Trace(ctx context.Context, begin time.Time, fc func() (command string, linesAffected int64), err error)
}
