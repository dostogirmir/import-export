package controllers

import (
	"github.com/gofiber/fiber/v2"
	"ispa-import-export/services"
	"log"
)

func ImportDistricts(c *fiber.Ctx) error {
	
	divisionID := c.FormValue("division_id")
    if divisionID == "" {
        return c.Status(fiber.StatusBadRequest).SendString("division_id is required")
    }
	  //Retrieve the uploaded file
	  file, err := c.FormFile("file")
	  if err != nil {
		  return c.Status(fiber.StatusBadRequest).SendString("File upload error: " + err.Error())
	  }
  
	  // Open the file
	  f, err := file.Open()
	  if err != nil {
		  return c.Status(fiber.StatusInternalServerError).SendString("Error opening file: " + err.Error())
	  }
	  defer f.Close()
  
	  // Call the service to process the CSV file
	  if err := services.ProcessDistrictCSV(f,divisionID); err != nil {
		  log.Printf("Error processing CSV: %v", err)
		  return c.Status(fiber.StatusInternalServerError).SendString("Error processing CSV file")
	  }
	
   

	return c.SendString("File imported successfully",)
}