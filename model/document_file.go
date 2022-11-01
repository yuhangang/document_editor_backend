package model

import (
	"time"

	"gorm.io/gorm"
)

type DocumentFile struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Title     string         `json:"title"`
	Data      string         `json:"data"`
	DeviceId  string         `gorm:"foreignkey:devic_info_id" json:"device_id"`
}
