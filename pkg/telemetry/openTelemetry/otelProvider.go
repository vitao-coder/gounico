package openTelemetry

import (
	"context"
	"gounico/global"
	"gounico/pkg/messaging/pulsar/tracing"

	"github.com/apache/pulsar-client-go/pulsar"
	"go.opentelemetry.io/otel/propagation"

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

func ExtractTraceContextFromMessage(ctx context.Context, message pulsar.ConsumerMessage) context.Context {
	consumerAdapter := tracing.ConsumerMessageAdapter{Message: message}
	traceContext := propagation.TraceContext{}
	return traceContext.Extract(ctx, &consumerAdapter)
}

func BuildAndInjectSpanOnMessageContext(ctx context.Context, injectedSpanName string, message *pulsar.ProducerMessage) context.Context {
	ctxWithSpan, traceSpan := NewSpan(ctx, injectedSpanName)
	traceContext := propagation.TraceContext{}
	message.Properties = map[string]string{}
	producerAdapter := tracing.ProducerMessageAdapter{Message: message}
	traceContext.Inject(ctxWithSpan, &producerAdapter)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		traceContext,
	))

	ctxInjected := ContextWithSpan(ctx, traceSpan)

	return ctxInjected
}

func buildAndInjectSpanOnContext(ctx context.Context, injectedSpanName string) context.Context {
	ctxWithSpan, traceSpan := NewSpan(ctx, injectedSpanName)
	traceContext := propagation.TraceContext{}
	contextAdapter := ContextAdapter{ctx}
	traceContext.Inject(ctxWithSpan, &contextAdapter)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		traceContext,
	))
	ctxInjected := ContextWithSpan(ctx, traceSpan)
	return ctxInjected
}

func extractTraceContextFromContext(ctx context.Context) context.Context {
	contextAdapter := ContextAdapter{ctx}
	traceContext := propagation.TraceContext{}
	return traceContext.Extract(ctx, &contextAdapter)
}

func TraceContextSpan(ctx context.Context, contextName string) (context.Context, trace.Span) {
	ctxExtracted := extractTraceContextFromContext(ctx)
	ctx = buildAndInjectSpanOnContext(ctxExtracted, contextName)
	traceSpan := SpanFromContext(ctx)
	return ctx, traceSpan
}
