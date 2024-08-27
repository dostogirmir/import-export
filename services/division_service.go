package services

import (
    "import-export/db"
    "import-export/models"
    "github.com/gofiber/fiber/v2"
    "encoding/csv"
    "log"
    "strconv"
    "sync"
    "bytes"
    "bufio"
	"gorm.io/gorm"
	"io"
)
type DivisionService struct {
    // Add necessary fields if needed
}

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


const (
    batchSize   = 20000 // Number of records per batch
    workerCount = 5    // Number of concurrent workers for batch insertion
)

// ProcessDivisionCSV processes the CSV file in a memory-efficient way
func ProcessDivisionCSV(reader io.Reader) error {
    db := db.GetDB() // Get the GORM database connection
    csvReader := csv.NewReader(bufio.NewReader(reader))

    // Skip the header line
    if _, err := csvReader.Read(); err != nil {
        return err
    }

    var wg sync.WaitGroup
    batchChan := make(chan []models.Division, workerCount) // Channel to hold batches

    // Start a fixed number of workers for inserting batches
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go worker(db, batchChan, &wg)
    }

    divisions := make([]models.Division, 0, batchSize)

    // Read and process the file line-by-line
    for {
        line, err := csvReader.Read()
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            log.Printf("Error reading line: %v", err)
            continue
        }

        // Parse CSV line to model
        division := models.Division{
            Name: line[1], 
            // Add other fields accordingly
        }

        // Add to the current batch
        divisions = append(divisions, division)

        // If batch size is reached, send it to the channel and reset
        if len(divisions) == batchSize {
            batchChan <- divisions
            divisions = make([]models.Division, 0, batchSize)
        }
    }

    // Send any remaining records in the last batch
    if len(divisions) > 0 {
        batchChan <- divisions
    }

    close(batchChan) // Close the channel to signal workers to stop
    wg.Wait()        // Wait for all workers to finish

    return nil
}

// worker function for processing batch inserts concurrently
func worker(db *gorm.DB, batchChan <-chan []models.Division, wg *sync.WaitGroup) {
    defer wg.Done()

    for divisions := range batchChan {
        if err := bulkInsertDivisions(db, divisions); err != nil {
            log.Printf("Error inserting batch: %v", err)
        }
    }
}

// bulkInsertDivisions inserts a batch of divisions into the database
func bulkInsertDivisions(db *gorm.DB, divisions []models.Division) error {
    tx := db.Begin()
    defer tx.Rollback()

    for i := 0; i < len(divisions); i += 1000 {
        chunk := divisions[i:min(i+1000, len(divisions))]
        if err := tx.CreateInBatches(&chunk, 1000).Error; err != nil {
            return err
        }
    }

    return tx.Commit().Error
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}