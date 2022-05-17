package openTelemetry

import (
	"context"

	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type spanWrapper struct {
	ctx  context.Context
	span trace.Span
}

func newSpanWrapper(ctx context.Context, span trace.Span) spanWrapper {
	return spanWrapper{
		ctx:  ctx,
		span: span,
	}
}

func (sw spanWrapper) AddSpanTags(tags map[string]string) {
	list := make([]attribute.KeyValue, len(tags))

	var i int
	for k, v := range tags {
		list[i] = attribute.Key(k).String(v)
		i++
	}

	sw.span.SetAttributes(list...)
}

func (sw spanWrapper) AddSpanEvents(name string, events map[string]string) {
	list := make([]trace.EventOption, len(events))

	var i int
	for k, v := range events {
		list[i] = trace.WithAttributes(attribute.Key(k).String(v))
		i++
	}

	sw.span.AddEvent(name, list...)
}

func (sw spanWrapper) AddSpanError(err error) {
	sw.span.RecordError(err)
}

func (sw spanWrapper) FailSpan(msg string) {
	sw.span.SetStatus(codes.Error, msg)
}

func (sw spanWrapper) SuccessSpan(msg string) {
	sw.span.SetStatus(codes.Ok, msg)
}

func (sw spanWrapper) End() {
	sw.span.End()
}
