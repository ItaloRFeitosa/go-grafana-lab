package order

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/italorfeitosa/go-grafana-lab/internal/order/model"
	warehouseclient "github.com/italorfeitosa/go-grafana-lab/internal/warehouse/client"
)

func SetRoutes(app *fiber.App) {
	warehouseClient := warehouseclient.New()

	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		order := model.Order{ID: c.Params("id")}

		ctx := c.UserContext()

		if err := callRepositoryOperation(ctx, "getOrder", order.ID); err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(order)
	})

	app.Patch("/orders/:id/approve", func(c *fiber.Ctx) error {
		order := model.Order{ID: c.Params("id")}

		ctx := c.UserContext()

		if err := callRepositoryOperation(ctx, "approveOrder", order.ID); err != nil {
			return err
		}

		if err := warehouseClient.OrderDispatch(ctx, order.ID); err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(order)
	})

	app.Patch("/orders/:id/fail", func(c *fiber.Ctx) error {
		order := model.Order{ID: c.Params("id")}

		ctx := c.UserContext()

		if err := callRepositoryOperation(ctx, "failOrder", order.ID); err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(order)
	})

	app.Post("/orders", func(c *fiber.Ctx) error {
		var order model.Order

		ctx := c.UserContext()

		if err := c.BodyParser(&order); err != nil {
			return err
		}

		if err := warehouseClient.OrderPrepare(ctx, order.ID); err != nil {
			return err
		}

		if err := callRepositoryOperation(ctx, "saveOrder", order.ID); err != nil {
			return err
		}

		return c.SendStatus(http.StatusAccepted)
	})

}
