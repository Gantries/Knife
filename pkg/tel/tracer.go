// Package tel contains opentelemetry helpers.
package tel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var noopTracerProvider = noop.NewTracerProvider()
var noopTracer = noopTracerProvider.Tracer("noop")
var tracer *trace.Tracer

func newTracer(name string) {
	t := otel.Tracer(name)
	tracer = &t
}

// SetupTracer sets the global tracer.
func SetupTracer(t *trace.Tracer) {
	tracer = t
}

// Span creates a new span.
func Span(ctx *context.Context, name string) (span trace.Span) {
	if tracer != nil {
		c, span := (*tracer).Start(*ctx, name)
		*ctx = c
		return span
	}
	*ctx, span = noopTracer.Start(*ctx, name)
	return
}

// Do ends the span.
func Do(s trace.Span, err *error) {
	defer s.End()
	var e error
	if err != nil {
		e = *err
	}
	if e == nil {
		s.SetStatus(codes.Ok, "")
	} else {
		s.SetStatus(codes.Error, e.Error())
	}
}
