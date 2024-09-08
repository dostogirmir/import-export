package main

import (
	"encoding/json"
	"log"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"ispa-import-export/db"
	"ispa-import-export/routes"
	"net/url"
    "github.com/gorilla/websocket"
)

type BroadcastMessage struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type CsvImportCompletedData struct {
	UserId  int    `json:"userId"`
	Message string `json:"message"`
}


func sendWebSocketNotification(userId int, message string) error {
	// Reverb server connection details
	u := url.URL{
		Scheme: "ws", // Use "wss" if using HTTPS
		Host:   "localhost:8080", // Change to your Reverb host and port
		Path:   fmt.Sprintf("/app/%s", "foqt9nb0m8xgluglxxso"), // Replace with your REVERB_APP_KEY
	}

	// Set headers for authentication if needed
	headers := map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", "28lzgkcxxowcpbhwuujo")}, // Optional if your server requires it
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}
	defer conn.Close()

	// Prepare the data for broadcasting
	data := CsvImportCompletedData{
		UserId:  userId,
		Message: message,
	}

	// Prepare the broadcast message
	broadcastMessage := BroadcastMessage{
		Event:   "CsvImportCompleted",
		Channel: fmt.Sprintf("private-import-notifications.%d", userId),
		Data:    data,
	}

	// Convert the broadcast message to JSON
	messageBytes, err := json.Marshal(broadcastMessage)
	if err != nil {
		return err
	}

	// Send the broadcast message over the WebSocket connection
	err = conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Initialize Fiber
	app := fiber.New(fiber.Config{
        // Set the maximum allowed size for uploads to 100MB
        BodyLimit: 100 * 1024 * 1024, // 100MB
    })

	
	// Initialize database connection
	err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Register all routes
	// routes.SetupCustomerRoutes(app)
	// routes.SetupAreaRoutes(app)
	routes.RegisterDivisionRoutes(app)
	routes.RegisterDistrictRoutes(app)
	// Create test callback route
   // Define a route that sends a message to the Laravel WebSocket server
    // Define a route that sends a message to the Laravel WebSocket server
	app.Get("/", func(c *fiber.Ctx) error {
		userId := 1
		message := "CSV import completed successfully!"

		err := sendWebSocketNotification(userId, message)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.SendString("Notification sent successfully!")
	})

	// Start server on port 3000
	log.Fatal(app.Listen(":3000"))
}