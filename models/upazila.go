package models

import (
    // "gorm.io/gorm"
)

type Upazila struct {
    ID         uint      `gorm:"primaryKey"`
    Name       string    `json:"name" gorm:"size:255;not null"`
    DistrictID uint      `json:"district_id" gorm:"not null"`
    District   District  `gorm:"foreignKey:DistrictID"`
}