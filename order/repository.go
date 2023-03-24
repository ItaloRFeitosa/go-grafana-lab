package order

import (
	"context"
	"fmt"

	"github.com/italorfeitosa/go-grafana-lab/chaos"
	"github.com/italorfeitosa/go-grafana-lab/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func orderOperation(ctx context.Context, operationName, id string) error {
	spanName := fmt.Sprint("repository::", operationName)
	attr := attribute.String("order.id", id)
	_, span := tracing.Start(ctx, spanName, trace.WithAttributes(attr))
	defer span.End()

	chaos.Latency()

	if err := chaos.Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
