package telemetry

import (
	"context"
)

type OpenTelemetry interface {
	Start(ctx context.Context, spanName string)
	AddSpanTags(tags map[string]string)
	AddSpanEvents(name string, events map[string]string)
	AddSpanError(err error)
	FailSpan(msg string)
	SuccessSpan(msg string)
	End()
}
