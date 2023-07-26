package checkout

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/italorfeitosa/go-grafana-lab/internal/checkout/model"
	orderclient "github.com/italorfeitosa/go-grafana-lab/internal/order/client"
	paymentclient "github.com/italorfeitosa/go-grafana-lab/internal/payment/client"
	paymentmodel "github.com/italorfeitosa/go-grafana-lab/internal/payment/model"
)

func SetRoutes(app *fiber.App) {
	orderClient := orderclient.New()
	paymentClient := paymentclient.New()

	app.Put("/checkouts/:id", func(c *fiber.Ctx) error {
		if err := orderClient.CreateOrder(c.UserContext(), c.Params("id")); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})

	app.Patch("/checkouts/:id/finish", func(c *fiber.Ctx) error {
		var (
			paym model.Payment
			err  error
		)

		id := c.Params("id")

		if err := c.BodyParser(&paym); err != nil {
			return err
		}

		if paym.Status == "PAID" {
			err = orderClient.ApproveOrder(c.UserContext(), id)
		} else {
			err = orderClient.FailOrder(c.UserContext(), id)
		}

		if err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})

	app.Post("/checkouts/:id/payments", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if _, err := orderClient.GetOrder(c.UserContext(), id); err != nil {
			return err
		}

		if err := paymentClient.CreatePayment(c.UserContext(), paymentmodel.Payment{
			ID:            uuid.NewString(),
			CorrelationID: id,
		}); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})
}
