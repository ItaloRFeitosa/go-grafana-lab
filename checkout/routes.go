package checkout

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/italorfeitosa/go-grafana-lab/order"
	"github.com/italorfeitosa/go-grafana-lab/payment"
)

func SetRoutes(app *fiber.App) {
	orderClient := order.NewClient()
	paymentClient := payment.NewClient()

	app.Put("/checkouts/:id", func(c *fiber.Ctx) error {
		if err := orderClient.CreateOrder(c.UserContext(), c.Params("id")); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})

	app.Patch("/checkouts/:id/finish", func(c *fiber.Ctx) error {
		var (
			paym Payment
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

		if err := paymentClient.CreatePayment(c.UserContext(), payment.Payment{
			ID:            uuid.NewString(),
			CorrelationID: id,
		}); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})
}
