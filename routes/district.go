package routes

import (
	"github.com/gofiber/fiber/v2"
	"ispa-import-export/controllers"
)

func RegisterDistrictRoutes(app *fiber.App) {
	districtGroup := app.Group("/districts")

	districtGroup.Post("/import", controllers.ImportDistricts)
	
}