package payment

import (
	"context"

	"github.com/italorfeitosa/go-grafana-lab/pkg/chaos"
	"github.com/italorfeitosa/go-grafana-lab/pkg/tracing"
	"go.opentelemetry.io/otel/codes"
)

func callPaymentGateway(ctx context.Context) error {
	_, span := tracing.Start(ctx, "gateway::paymentAttempt")
	defer span.End()

	chaos.Latency()

	if err := chaos.Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
