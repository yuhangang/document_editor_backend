package model

import (
	"gorm.io/gorm"
)

type Continent struct {
	Code string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type RecordContinent struct {
	Code string
	Name string
}

func MigrateContinent(db *gorm.DB) error {
	err := db.AutoMigrate(&Continent{})
	return err
}