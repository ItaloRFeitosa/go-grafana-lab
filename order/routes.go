package order

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/italorfeitosa/go-grafana-lab/warehouse"
)

func SetRoutes(app *fiber.App) {
	warehouseClient := warehouse.NewClient()

	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		order := Order{ID: c.Params("id")}

		ctx := c.UserContext()

		if err := orderOperation(ctx, "getOrder", order.ID); err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(order)
	})

	app.Patch("/orders/:id/approve", func(c *fiber.Ctx) error {
		order := Order{ID: c.Params("id")}

		ctx := c.UserContext()

		if err := orderOperation(ctx, "approveOrder", order.ID); err != nil {
			return err
		}

		if err := warehouseClient.OrderDispatch(ctx, order.ID); err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(order)
	})

	app.Patch("/orders/:id/fail", func(c *fiber.Ctx) error {
		order := Order{ID: c.Params("id")}

		ctx := c.UserContext()

		if err := orderOperation(ctx, "failOrder", order.ID); err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(order)
	})

	app.Post("/orders", func(c *fiber.Ctx) error {
		var order Order

		ctx := c.UserContext()

		if err := c.BodyParser(&order); err != nil {
			return err
		}

		if err := warehouseClient.OrderPrepare(ctx, order.ID); err != nil {
			return err
		}

		if err := orderOperation(ctx, "saveOrder", order.ID); err != nil {
			return err
		}

		return c.SendStatus(http.StatusAccepted)
	})

}
