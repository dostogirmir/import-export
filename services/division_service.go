package services

import (
	"import-export/db"
	"import-export/models"
	"encoding/csv"
	"sync"
	"bytes"
	"log"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"bufio"
)

type DivisionService struct{}

func NewDivisionService() *DivisionService {
	return &DivisionService{}
}

// GetAllDivisions returns all divisions
func GetAllDivisions(page int, pageSize int) ([]models.Division, error) {
	var divisions []models.Division
	result := db.GetDB().Limit(pageSize).Offset((page - 1) * pageSize).Find(&divisions)
	if result.Error != nil {
		return nil, result.Error
	}
	return divisions, nil
}

// GetDivisionByID returns a division by its ID
func GetDivisionByID(id uint) (*models.Division, error) {
	var division models.Division
	result := db.GetDB().First(&division, id) // Use GetDB() to access the DB instance
	if result.Error != nil {
		return nil, result.Error
	}
	return &division, nil
}

// CreateDivision creates a new division
func CreateDivision(division *models.Division) error {
	result := db.GetDB().Create(division) // Use GetDB() to access the DB instance
	return result.Error
}

// UpdateDivision updates an existing division
func UpdateDivision(id uint, division *models.Division) error {
	var existingDivision models.Division
	result := db.GetDB().First(&existingDivision, id) // Use GetDB() to access the DB instance
	if result.Error != nil {
		return result.Error
	}
	existingDivision.Name = division.Name
	db.GetDB().Save(&existingDivision) // Use GetDB() to access the DB instance
	return nil
}

// DeleteDivision deletes a division by its ID
func DeleteDivision(id uint) error {
	var division models.Division
	result := db.GetDB().First(&division, id) // Use GetDB() to access the DB instance
	if result.Error != nil {
		return result.Error
	}
	db.GetDB().Delete(&division) // Use GetDB() to access the DB instance
	return nil
}

func (s *DivisionService) ExportDivisions(c *fiber.Ctx) error {
	chunkSize := 10000 // adjust the chunk size based on your system's performance
	offset := 0

	// create a channel to receive the exported data
	exportChan := make(chan []byte, chunkSize)

	// create a wait group to wait for all goroutines to finish
	wg := &sync.WaitGroup{}

	// fetch and process the data in chunks
	for {
		divisions, err := s.getDivisionsChunk(chunkSize, offset) // fetch the next chunk of divisions
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve divisions"})
		}

		if len(divisions) == 0 {
			break // no more data, exit the loop
		}

		wg.Add(1)
		go func(divisions []models.Division) {
			defer wg.Done()
			exportData, err := s.exportDivisionsChunk(divisions)
			if err != nil {
				log.Println(err)
				return
			}
			exportChan <- exportData
		}(divisions)

		offset += chunkSize
	}

	// wait for all goroutines to finish
	wg.Wait()

	// close the channel
	close(exportChan)

	// set the response headers
	c.Response().Header.Set("Content-Type", "text/csv")
	c.Response().Header.Set("Content-Disposition", "attachment; filename=\"divisions.csv\"")

	// create a CSV writer
	writer := bufio.NewWriter(c.Response().BodyWriter())

	csvWriter := csv.NewWriter(writer)

	// write the header row
	header := []string{"ID", "Name"}
	if err := csvWriter.Write(header); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to write CSV header"})
	}

	// write each chunk of data to the CSV file
	for exportData := range exportChan {
		writer.Write(exportData)
	}

	// flush the CSV writer
	csvWriter.Flush()
	writer.Flush()

	return nil
}

func (s *DivisionService) getDivisionsChunk(chunkSize int, offset int) ([]models.Division, error) {
	var divisions []models.Division

	// Use GORM to fetch the next chunk of divisions
	err := db.GetDB().Model(&models.Division{}).
		Limit(chunkSize).
		Offset(offset).
		Find(&divisions).
		Error

	if err != nil {
		return nil, err
	}

	return divisions, nil
}

func (s *DivisionService) exportDivisionsChunk(divisions []models.Division) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// write each division as a row
	for _, division := range divisions {
		row := []string{strconv.Itoa(int(division.ID)), division.Name}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	// flush the writer
	writer.Flush()

	return buf.Bytes(), nil
}