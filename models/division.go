package models

import (
    "time"
)

type Division struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:255;not null;unique"`
    CreatedAt time.Time      `gorm:"autoCreateTime"` // Automatically set on creation
    UpdatedAt time.Time      `gorm:"autoUpdateTime"` // Automatically updated on each update
}
