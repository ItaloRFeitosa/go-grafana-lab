package payment

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	checkoutclient "github.com/italorfeitosa/go-grafana-lab/internal/checkout/client"
	checkoutmodel "github.com/italorfeitosa/go-grafana-lab/internal/checkout/model"
	"github.com/italorfeitosa/go-grafana-lab/internal/payment/model"
	"github.com/italorfeitosa/go-grafana-lab/pkg/chaos"
	"github.com/italorfeitosa/go-grafana-lab/pkg/tracing"
)

var checkoutClient = checkoutclient.New()

func SetRoutes(app *fiber.App) {
	app.Post("/payments", func(c *fiber.Ctx) error {
		var payment model.Payment
		payment.Status = model.StatusPending

		ctx := c.UserContext()

		if err := c.BodyParser(&payment); err != nil {
			return err
		}

		if err := callRepositoryOperation(ctx, "savePayment", payment); err != nil {
			return err
		}

		gatewayError := callPaymentGateway(ctx)
		if gatewayError == nil {
			go simulateWebhookCall(payment)

			return c.SendStatus(http.StatusAccepted)
		}

		payment.Status = model.StatusFailed

		if err := callRepositoryOperation(ctx, "savePayment", payment); err != nil {
			return err
		}

		return gatewayError
	})

}

func simulateWebhookCall(p model.Payment) {
	chaos.Latency()
	ctx := context.Background()
	ctx, span := tracing.Start(ctx, "webhook::paymentConfirmation")
	defer span.End()

	p.Status = "PAID"
	if err := chaos.Error(); err != nil {
		p.Status = "FAILED"
	}

	callRepositoryOperation(ctx, "savePayment", p)

	checkoutClient.FinishCheckout(ctx, checkoutmodel.Payment{
		CorrelationID: p.CorrelationID,
		Status:        p.Status,
	})
}
