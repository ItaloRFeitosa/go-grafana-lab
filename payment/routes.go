package payment

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/italorfeitosa/go-grafana-lab/chaos"
	"github.com/italorfeitosa/go-grafana-lab/tracing"
	"github.com/italorfeitosa/go-grafana-lab/webhook"
)

func SetRoutes(app *fiber.App) {
	webhookClient := webhook.NewCheckoutClient()
	app.Post("/payments", func(c *fiber.Ctx) error {
		var payment Payment
		payment.Status = StatusPending

		ctx := c.UserContext()

		if err := c.BodyParser(&payment); err != nil {
			return err
		}

		if err := paymentOperation(ctx, "savePayment", payment); err != nil {
			return err
		}

		gatewayError := CallPaymentGateway(ctx)
		if gatewayError == nil {
			go func(p Payment) {
				ctx := context.Background()
				time.Sleep(3 * time.Second)
				ctx, span := tracing.Start(ctx, "webhook::paymentConfirmation")
				defer span.End()

				p.Status = "PAID"
				if err := chaos.Error(); err != nil {
					p.Status = "FAILED"
				}

				paymentOperation(ctx, "savePayment", payment)

				webhookClient.FinishCheckout(ctx, webhook.Payment{
					CorrelationID: p.CorrelationID,
					Status:        p.Status,
				})

			}(payment)

			return c.SendStatus(http.StatusAccepted)
		}

		payment.Status = StatusFailed

		if err := paymentOperation(ctx, "savePayment", payment); err != nil {
			return err
		}

		return gatewayError
	})

}
