package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func setTracer(t trace.Tracer) {
	tracer = t
}

func Tracer() trace.Tracer {
	return tracer
}

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, opts...)
}
