package model

import (
	"time"

	"gorm.io/gorm"
)

type DeviceInfo struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DeviceId              string         `json:"device_id" gorm:"primaryKey"`
	DeviceModel           string         `json:"device_model"`
	DeviceManufacturer    string         `json:"device_manufacturer"`
	DeviceOsVersion       string         `json:"device_os_version"`
	DeviceOsVersionNumber *float64       `json:"device_os_version_number"`
	Lat                   *float64       `json:"lat"`
	Lng                   *float64       `json:"lng"`
	DeviceUserId          *uint          `gorm:"foreignkey:country_id" json:"device_user_id"`
}
