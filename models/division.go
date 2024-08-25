package models

type Division struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `json:"name" gorm:"size:255;not null"`
}
