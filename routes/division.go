// routes/division.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"import-export/controllers"
)

func RegisterDivisionRoutes(app *fiber.App) {
	divisionGroup := app.Group("/divisions")

	divisionGroup.Get("/", controllers.GetAllDivisions)
	divisionGroup.Get("/:id", controllers.GetDivisionByID)
	divisionGroup.Post("/", controllers.CreateDivision)
	divisionGroup.Put("/:id", controllers.UpdateDivision)
	divisionGroup.Delete("/:id", controllers.DeleteDivision)
    divisionGroup.Post("/import", controllers.ImportDivisions)
    divisionGroup.Get("/export", controllers.ExportDivisions)
}
