package services

import (
    "ispa-import-export/db"
    "ispa-import-export/models"
    "github.com/gofiber/fiber/v2"
    "encoding/csv"
    "log"
    // "strconv"
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

// BulkDeleteDivisionsResponse deletes multiple divisions by their IDs in batches and returns a response.
func BulkDeleteDivisionsResponse(ids []int) (map[string]interface{}, error) {
    batchSize := 10000             // Batch size to process IDs in chunks
    numWorkers := 5                // Number of concurrent workers

    var wg sync.WaitGroup
    sem := make(chan struct{}, numWorkers) // Semaphore to limit concurrent workers
    deleteCount := 0 // Count of deleted records
    var mu sync.Mutex // Mutex to protect shared resources

    for i := 0; i < len(ids); i += batchSize {
        end := i + batchSize
        if end > len(ids) {
            end = len(ids)
        }
        batchIds := ids[i:end]

        wg.Add(1)
        sem <- struct{}{}

        go func(batchIds []int) {
            defer wg.Done()
            defer func() { <-sem }()

            if err := deleteBatch(batchIds); err != nil {
                log.Printf("Error deleting batch: %v\n", err)
                return
            }

            // Update the delete count
            mu.Lock()
            deleteCount += len(batchIds)
            mu.Unlock()

        }(batchIds)
    }

    wg.Wait()

    // Prepare a response
    response := map[string]interface{}{
        "deleted_records": deleteCount,
    }

    return response, nil
}

// deleteBatch deletes a batch of divisions by their IDs.
func deleteBatch(ids []int) error {
    tx := db.GetDB().Begin() // Start a new transaction

    // Use a single delete statement with IN clause for efficiency
    if err := tx.Where("id IN ?", ids).Delete(&models.Division{}).Error; err != nil {
        tx.Rollback() // Rollback in case of error
        return err
    }

    return tx.Commit().Error // Commit the transaction
}

func (s *DivisionService) ExportDivisions(c *fiber.Ctx) error {
	chunkSize := 20000 // adjust the chunk size based on your system's performance
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
	header := []string{"Name"}
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
		row := []string{division.Name}
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
    db, csvReader := db.GetDB(), csv.NewReader(bufio.NewReader(reader))

    if _, err := csvReader.Read(); err != nil {
        return err
    }

    batchChan := make(chan []models.Division, workerCount)
    var wg sync.WaitGroup

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for divisions := range batchChan {
                if err := bulkInsertDivisions(db, divisions); err != nil {
                    log.Printf("Batch insert error: %v", err)
                }
            }
        }()
    }

    divisions := make([]models.Division, 0, batchSize)
    for line, err := csvReader.Read(); err == nil; line, err = csvReader.Read() {
        divisions = append(divisions, models.Division{
            Name: line[0],
        })
        if len(divisions) == batchSize {
            batchChan <- divisions
            divisions = make([]models.Division, 0, batchSize)
        }
    }

    if len(divisions) > 0 {
        batchChan <- divisions
    }
    
    close(batchChan)
    wg.Wait()

    return nil
}

// bulkInsertDivisions inserts a batch of divisions into the database
func bulkInsertDivisions(db *gorm.DB, divisions []models.Division) error {
    for i := 0; i < len(divisions); i += 1000 {
        end := i + 1000
        if end > len(divisions) {
            end = len(divisions) // Ensure we don't go beyond the slice bounds
        }
        chunk := divisions[i:end]
        if err := db.CreateInBatches(&chunk, 1000).Error; err != nil {
            return err
        }
    }
    return nil
}
