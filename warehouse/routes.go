package warehouse

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	app.Patch("/warehouse/orders/:id/dispatch", func(c *fiber.Ctx) error {
		if err := SendMessageToBroker(c.UserContext()); err != nil {
			return err
		}

		if err := orderOperation(c.UserContext(), "markAsAwaitingDispatch", c.Params("id")); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})

	app.Patch("/warehouse/orders/:id/prepare", func(c *fiber.Ctx) error {
		if err := orderOperation(c.UserContext(), "saveOrder", c.Params("id")); err != nil {
			return err
		}

		if err := orderOperation(c.UserContext(), "reserveOrderItems", c.Params("id")); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})

}
