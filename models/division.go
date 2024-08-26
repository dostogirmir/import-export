package models

import (
    "time"
)

type Division struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"size:255;not null;unique"` // Add unique constraint
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
