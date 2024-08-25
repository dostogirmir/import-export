package models

import (
    // "gorm.io/gorm"
)

type District struct {
    ID         uint      `gorm:"primaryKey"`
    Name       string    `json:"name" gorm:"size:255;not null"`
    DivisionID uint      `json:"division_id" gorm:"not null"`
    Division   Division  `gorm:"foreignKey:DivisionID"`
}