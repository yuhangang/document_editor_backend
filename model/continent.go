package model

import (
	"gorm.io/gorm"
)

type Continent struct {
	Code      string    `gorm:"primaryKey" json:"code"`
	Name      string    `gorm:"unique" json:"name"`
	Countries []Country `json:"countries"`
}

type RecordContinent struct {
	Code string
	Name string
}

func MigrateContinent(db *gorm.DB) error {
	err := db.AutoMigrate(&Continent{})
	return err
}
