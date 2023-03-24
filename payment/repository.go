package payment

import (
	"context"
	"fmt"

	"github.com/italorfeitosa/go-grafana-lab/chaos"
	"github.com/italorfeitosa/go-grafana-lab/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func paymentOperation(ctx context.Context, operationName string, payment Payment) error {
	spanName := fmt.Sprint("repository::", operationName)
	attrs := []attribute.KeyValue{
		attribute.String("payment.correlation.id", payment.CorrelationID),
		attribute.String("payment.id", payment.ID),
		attribute.String("payment.status", payment.Status),
	}
	_, span := tracing.Start(ctx, spanName, trace.WithAttributes(attrs...))

	defer span.End()

	chaos.Latency()

	if err := chaos.Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
