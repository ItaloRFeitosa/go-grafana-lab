package payment

import (
	"context"

	"github.com/italorfeitosa/go-grafana-lab/chaos"
	"github.com/italorfeitosa/go-grafana-lab/tracing"
	"go.opentelemetry.io/otel/codes"
)

func CallPaymentGateway(ctx context.Context) error {
	_, span := tracing.Start(ctx, "gateway::paymentAttempt")
	defer span.End()

	chaos.Latency()

	if err := chaos.Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
