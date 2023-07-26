package payment

import (
	"context"
	"fmt"

	"github.com/italorfeitosa/go-grafana-lab/internal/payment/model"
	"github.com/italorfeitosa/go-grafana-lab/pkg/chaos"
	"github.com/italorfeitosa/go-grafana-lab/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func callRepositoryOperation(ctx context.Context, operationName string, payment model.Payment) error {
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
