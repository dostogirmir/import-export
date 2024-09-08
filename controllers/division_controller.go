// controllers/division_controller.go
package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"ispa-import-export/models"
	"ispa-import-export/services"
	"log"
	"strings" // Make sure this import is included
)

// Get all divisions
func GetAllDivisions(c *fiber.Ctx) error {
	divisions, err := services.GetAllDivisions( 1, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve divisions"})
	}
	return c.JSON(divisions)
}

// Get division by ID
func GetDivisionByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid division ID"})
	}

	division, err := services.GetDivisionByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Division not found"})
	}
	return c.JSON(division)
}

// Create a new division
func CreateDivision(c *fiber.Ctx) error {
	division := new(models.Division)
	if err := c.BodyParser(division); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	if err := services.CreateDivision(division); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create division"})
	}

	return c.Status(fiber.StatusCreated).JSON(division)
}

// Update an existing division
func UpdateDivision(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid division ID"})
	}

	division := new(models.Division)
	if err := c.BodyParser(division); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	if err := services.UpdateDivision(uint(id), division); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update division"})
	}

	return c.JSON(fiber.Map{"message": "Division updated successfully"})
}

// Delete a division
func DeleteDivision(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid division ID"})
	}

	if err := services.DeleteDivision(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete division"})
	}

	return c.JSON(fiber.Map{"message": "Division deleted successfully"})
}

func BulkDeleteDivisions(c *fiber.Ctx) error {
    rawBody := c.Body()
    log.Println("Raw body received:", string(rawBody))

    var ids []int
    if err := c.BodyParser(&ids); err != nil {
        log.Println("Error parsing body:", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    log.Println("BulkDeleteDivisions IDs:", ids)
    response, err := services.BulkDeleteDivisionsResponse(ids)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete divisions"})
    }

    return c.Status(fiber.StatusOK).JSON(response)
}




// ExportDivisions exports all divisions to a CSV file
func ExportDivisions(c *fiber.Ctx) error {
	divisionsService := services.NewDivisionService()
	return divisionsService.ExportDivisions(c)
}

// sanitizeFilename sanitizes file names to avoid issues with file paths
func sanitizeFilename(filename string) string {
    // Replace invalid characters and avoid directory traversal
    return strings.ReplaceAll(filename, "/", "_")
}

// ImportDivisions handles the CSV import request
func ImportDivisions(c *fiber.Ctx) error {
    // Retrieve the uploaded file
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
    if err := services.ProcessDivisionCSV(f); err != nil {
        log.Printf("Error processing CSV: %v", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Error processing CSV file")
    }

    return c.SendString("File imported successfully")
}