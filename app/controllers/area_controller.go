package controllers

import (
    "encoding/csv"
    "github.com/gofiber/fiber/v2"
    "github.com/xuri/excelize/v2"
    "io"
    "import-export/app/models"
    "import-export/app/services"
    "strings"
)

func BulkInsertArea(c *fiber.Ctx) error {
    // Check if a file is attached
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
    }

    // Open the file
    src, err := file.Open()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
    }
    defer src.Close()

    var areas []models.Area

    // Determine file type and process accordingly
    if strings.HasSuffix(file.Filename, ".csv") {
        reader := csv.NewReader(src)
        records, err := reader.ReadAll()
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse CSV"})
        }

        for _, record := range records {
            areas = append(areas, models.Area{Name: record[0]})
        }
    } else if strings.HasSuffix(file.Filename, ".xlsx") {
        xlsx, err := excelize.OpenReader(src)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse Excel"})
        }
        sheet := xlsx.GetSheetName(1)
        rows, err := xlsx.GetRows(sheet)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get rows"})
        }

        for _, row := range rows {
            if len(row) > 0 {
                areas = append(areas, models.Area{Name: row[0]})
            }
        }
    } else {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported file type"})
    }

    // Insert areas into the database
    err = services.BulkInsertAreas(areas)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert data"})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Data inserted successfully"})
}
