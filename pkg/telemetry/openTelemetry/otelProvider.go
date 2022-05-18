package openTelemetry

import (
	"context"
	"gounico/global"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Provider struct {
	provider trace.TracerProvider
	appName  string
}

func NewTracer(url string, appName string) *Provider {
	tp, err := jaegerProvider(url, appName)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)

	otc := &Provider{
		provider: tp,
		appName:  appName,
	}
	global.GlobalTracer = otel.Tracer(global.AppName)
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

func (p *Provider) Close(ctx context.Context) error {
	if prv, ok := p.provider.(*tracesdk.TracerProvider); ok {
		return prv.Shutdown(ctx)
	}
	return nil
}
