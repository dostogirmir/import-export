package models

import (
    "time"
)

type District struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:255;not null;unique"`
    DivisionID uint      `json:"division_id" gorm:"not null"`
    Division   Division  `gorm:"foreignKey:DivisionID"`
    CreatedAt time.Time      `gorm:"autoCreateTime"` 
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
}