package telemetry

import (
	"context"
)

type OpenTelemetry interface {
	Close(ctx context.Context) error
}
