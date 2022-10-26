package model

import "gorm.io/gorm"

type DeviceInfo struct {
	gorm.Model
	DeviceId              string   `json:"device_id" gorm:"unique"`
	DeviceModel           string   `json:"device_model"`
	DeviceManufacturer    string   `json:"device_manufacturer"`
	DeviceOsVersion       string   `json:"device_os_version"`
	DeviceOsVersionNumber string   `json:"device_os_version_number"`
	Lat                   *float64 `json:"lat"`
	Lng                   *float64 `json:"lng"`
	DeviceUserId          *uint    `json:"device_user_id"`
}
