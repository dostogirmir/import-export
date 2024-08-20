package routes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupDivisionRoutes sets up routes related to division functionality.
func SetupDivisionRoutes(app *fiber.App) {
	app.Get("/divisions", getDivisions)
	app.Post("/divisions", createDivision)
	// Add other division routes here
}

func getDivisions(c *fiber.Ctx) error {
	return c.SendString("List of divisions")
}

func createDivision(c *fiber.Ctx) error {
	return c.SendString("Create a division")
}
