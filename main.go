package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
 	"import-export/db"
	"import-export/routes"

)
func main() {
	// Initialize Fiber
	app := fiber.New()

	 // Initialize database connection
	// if err := db.InitDatabase(); err != nil {
    //     panic("Failed to connect to database")
    // }
	// Close the database connection when the main function is done

	// Register all routes
	routes.SetupCustomerRoutes(app)
	routes.SetupAreaRoutes(app)
	routes.SetupDivisionRoutes(app)

	// Start server on port 3000
	log.Fatal(app.Listen(":3000"))
}