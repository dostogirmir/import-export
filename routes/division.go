// routes/division.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"ispa-import-export/controllers"
)

func RegisterDivisionRoutes(app *fiber.App) {
	divisionGroup := app.Group("/divisions")

	divisionGroup.Get("/", controllers.GetAllDivisions)
	divisionGroup.Get("/:id", controllers.GetDivisionByID)
	divisionGroup.Post("/", controllers.CreateDivision)
	divisionGroup.Put("/:id", controllers.UpdateDivision)
	divisionGroup.Delete("/bulk-delete", controllers.BulkDeleteDivisions)
	divisionGroup.Delete("/:id", controllers.DeleteDivision)
    divisionGroup.Post("/import", controllers.ImportDivisions)
    divisionGroup.Get("/export/division", controllers.ExportDivisions)
}
