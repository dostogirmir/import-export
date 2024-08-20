package routes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupCustomerRoutes sets up routes related to customer functionality.
func SetupCustomerRoutes(app *fiber.App) {
	app.Get("/customers", getCustomers)
	app.Post("/customers", createCustomer)
	// Add other customer routes here
}

func getCustomers(c *fiber.Ctx) error {
	return c.SendString("List of customers")
}

func createCustomer(c *fiber.Ctx) error {
	return c.SendString("Create a customer")
}
