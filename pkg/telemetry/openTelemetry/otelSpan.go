package openTelemetry

import (
	"context"
	"gounico/global"

	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return global.GlobalTracer.Start(ctx, spanName)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func ContextWithSpan(parentContext context.Context, span trace.Span) context.Context {
	return trace.ContextWithSpan(parentContext, span)
}

func ContextWithSpanContext(parentContext context.Context, span trace.Span) context.Context {
	return trace.ContextWithSpanContext(parentContext, span.SpanContext())
}

func ContextWithRemoteSpanContext(parentContext context.Context, span trace.Span) context.Context {
	return trace.ContextWithRemoteSpanContext(parentContext, span.SpanContext())
}

func AddSpanTags(span trace.Span, tags map[string]string) {
	list := make([]attribute.KeyValue, len(tags))

	var i int
	for k, v := range tags {
		list[i] = attribute.Key(k).String(v)
		i++
	}

	span.SetAttributes(list...)
}

func AddSpanEvents(span trace.Span, name string, events map[string]string) {
	list := make([]trace.EventOption, len(events))

	var i int
	for k, v := range events {
		list[i] = trace.WithAttributes(attribute.Key(k).String(v))
		i++
	}

	span.AddEvent(name, list...)
}

func AddSpanError(span trace.Span, err error) {
	span.RecordError(err)
}

func FailSpan(span trace.Span, msg string) {
	span.SetStatus(codes.Error, msg)
}

func SuccessSpan(span trace.Span, msg string) {
	span.SetStatus(codes.Ok, msg)
}
