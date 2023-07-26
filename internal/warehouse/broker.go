package warehouse

import (
	"context"

	"github.com/italorfeitosa/go-grafana-lab/pkg/chaos"
	"github.com/italorfeitosa/go-grafana-lab/pkg/tracing"
	"go.opentelemetry.io/otel/codes"
)

func sendMessageToBroker(ctx context.Context) error {
	ctx, span := tracing.Start(ctx, "broker::sendMessage")
	defer span.End()

	chaos.Latency()

	if err := chaos.Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
