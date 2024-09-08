package services

import (
	"io"
	"log"
	"sync"
	"encoding/csv"
	"bufio"
	"strconv"
	"gorm.io/gorm"
	"ispa-import-export/db"
	"ispa-import-export/models"
)

// ProcessDistrictCSV processes the CSV file in a memory-efficient way
func ProcessDistrictCSV(reader io.Reader, divisionID string) error {
	db, csvReader := db.GetDB(), csv.NewReader(bufio.NewReader(reader))

    if _, err := csvReader.Read(); err != nil {
        return err
    }

	// Convert string to uint
	divID, err := strconv.ParseUint(divisionID, 10, 32)
	if err != nil {
		return err
	}
	divisionIDUint := uint(divID)

    batchChan := make(chan []models.District, workerCount)
    var wg sync.WaitGroup

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for districts := range batchChan {
                if err := bulkInsertDistricts(db, districts); err != nil {
                    log.Printf("Batch insert error: %v", err)
                }
            }
        }()
    }

    districts := make([]models.District, 0, batchSize)
    for line, err := csvReader.Read(); err == nil; line, err = csvReader.Read() {
        districts = append(districts, models.District{
            Name: line[0],
			DivisionID: divisionIDUint,
        })
        if len(districts) == batchSize {
            batchChan <- districts
            districts = make([]models.District, 0, batchSize)
        }
    }

    if len(districts) > 0 {
        batchChan <- districts
    }
    
    close(batchChan)
    wg.Wait()

    return nil
}

// bulkInsertdistricts inserts a batch of districts into the database
func bulkInsertDistricts(db *gorm.DB, districts []models.District) error {
    for i := 0; i < len(districts); i += 1000 {
        end := i + 1000
        if end > len(districts) {
            end = len(districts) // Ensure we don't go beyond the slice bounds
        }
        chunk := districts[i:end]
        if err := db.CreateInBatches(&chunk, 1000).Error; err != nil {
            return err
        }
    }
    return nil
}