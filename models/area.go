package models

import (
    // "gorm.io/gorm"
)

type Area struct {
    ID         uint   `gorm:"primaryKey"`
    Name       string `json:"name" gorm:"size:255;not null"`
    DivisionID uint   `json:"division_id" gorm:"not null"`
    DistrictID uint   `json:"district_id" gorm:"not null"`
    UpazilaID  uint   `json:"upazila_id" gorm:"not null"`
    CreatedBy  uint   `json:"created_by" gorm:"not null"`
}