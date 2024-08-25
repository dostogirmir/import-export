// controllers/division_controller.go
package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"import-export/models"
	"import-export/services"
	"os"
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

func ImportDivisions(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        log.Println("Error getting file:", err)
        return c.Status(fiber.StatusBadRequest).SendString("Failed to get file")
    }

    tempDir := "../temp/"
    if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
        log.Println("Error creating temp directory:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to create temp directory")
    }

    tempFile := tempDir + sanitizeFilename(file.Filename)
    log.Println("Saving file to:", tempFile)
    if err := c.SaveFile(file, tempFile); err != nil {
        log.Println("Error saving file:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
    }

    divisionsService := services.NewDivisionService()
    err = divisionsService.ImportDivisions(tempFile)
    if err != nil {
        log.Println("Error during import:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to import divisions")
    }

    if err := os.Remove(tempFile); err != nil {
        log.Println("Error removing temp file:", err)
    }

    return c.SendString("Divisions imported successfully")
}
