package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"import-export/db"
	"import-export/routes"
)

func main() {
	// Initialize Fiber
	app := fiber.New()

	
	// Initialize database connection
	err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Register all routes
	// routes.SetupCustomerRoutes(app)
	// routes.SetupAreaRoutes(app)
	routes.RegisterDivisionRoutes(app)

	// Start server on port 3000
	log.Fatal(app.Listen(":3000"))
}