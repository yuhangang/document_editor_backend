package model

import (
	"time"

	"gorm.io/gorm"
)

type DocumentFile struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `json:"title"`
	Data      string         `json:"data"`
}
