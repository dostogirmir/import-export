package services

import (
    "context"
    "import-export/db"
    "import-export/app/models"
    "gorm.io/gorm"
    "log"
    "time"
)

const chunkSize = 1000

func BulkInsertAreas(areas []models.Area) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()

    // Prepare the database connection
    tx := db.DB.Begin()
    if err := tx.Error; err != nil {
        return err
    }

    // Chunk the areas slice
    for i := 0; i < len(areas); i += chunkSize {
        end := i + chunkSize
        if end > len(areas) {
            end = len(areas)
        }
        chunk := areas[i:end]

        // Perform bulk insert in a separate goroutine
        go func(chunk []models.Area) {
            tx := db.DB.Begin()
            if err := tx.Error; err != nil {
                log.Println("Transaction begin error:", err)
                return
            }

            if err := tx.Create(&chunk).Error; err != nil {
                tx.Rollback()
                log.Println("Bulk insert error:", err)
                return
            }

            if err := tx.Commit().Error; err != nil {
                log.Println("Transaction commit error:", err)
            }
        }(chunk)
    }

    // Wait for all goroutines to complete
    time.Sleep(10 * time.Minute) // Adjust based on expected duration

    return nil
}
