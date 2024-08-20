package models

import (
    "gorm.io/gorm"
)

type Area struct {
    gorm.Model
    Name string `json:"name" gorm:"size:255;not null"`
}
