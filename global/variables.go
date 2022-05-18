package global

import (
	"go.opentelemetry.io/otel/trace"
)

var AppName = ""

var GlobalTracer trace.Tracer
