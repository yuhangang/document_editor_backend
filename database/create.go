package database

import (
	"echoapp/model"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
	TimeZone string
}

func CreateDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		" password=%s dbname=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	//db.AutoMigrate(&model.Continent{})
	//db.AutoMigrate(&model.Country{})

	db.Migrator().DropTable(&model.Continent{})
	db.Migrator().DropTable(&model.Country{})
	db.Migrator().DropTable(&model.City{})
	db.Migrator().CreateTable(&model.Continent{})
	db.Migrator().CreateTable(&model.Country{})
	db.Migrator().CreateTable(&model.City{})

	// Create
	db.Create(&model.Continent{Code: "AS", Name: "Asia"})
	db.Create(&model.Country{Code: "MY", Name: "Malaysia", Capital: "Kuala Lumpur", Native: "Malay", ContinentID: "AS"})

	// Read
	var continent model.Continent
	db.First(&continent, "code = ?", "AS") // find product with code D42

	// Update - update multiple fields
	db.Model(&continent).Updates(model.Continent{Code: "AS", Name: "Asia new"})
	return db, nil
}
