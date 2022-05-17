package openTelemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type openTelemetryTracer struct {
	otelTracer trace.Tracer
	spanWrapper
}

func NewTracer(url string, appName string) *openTelemetryTracer {
	tp, err := jaegerProvider(url, appName)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	otc := &openTelemetryTracer{otelTracer: otel.Tracer(appName)}
	return otc
}

func jaegerProvider(url string, appName string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		)),
	)
	return tp, nil
}

func (ott *openTelemetryTracer) Start(ctx context.Context, spanName string) {
	context, span := ott.otelTracer.Start(ctx, spanName)
	ott.spanWrapper = newSpanWrapper(context, span)
	return
}
