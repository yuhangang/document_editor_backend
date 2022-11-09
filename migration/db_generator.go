package migration

import (
	"echoapp/config"
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

func CreateDB(appConfig *config.CommandArgs) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		" password=%s dbname=%s sslmode=%s TimeZone=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASS"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_SSLMODE"),
		os.Getenv("POSTGRES_TIMEZONE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	//db.AutoMigrate(&model.Continent{})
	//db.AutoMigrate(&model.Country{})
	if appConfig.ClearDB {
		db.Migrator().AutoMigrate(&model.Continent{})
		db.Migrator().AutoMigrate(&model.Country{})
		db.Migrator().AutoMigrate(&model.City{})
		db.Migrator().AutoMigrate(&model.DeviceInfo{})
		db.Migrator().AutoMigrate(&model.DocumentFile{})
	} else {

		db.Migrator().DropTable(&model.Continent{})
		db.Migrator().DropTable(&model.Country{})
		db.Migrator().DropTable(&model.City{})
		db.Migrator().DropTable(&model.DeviceInfo{})
		db.Migrator().DropTable(&model.DocumentFile{})
		db.Migrator().CreateTable(&model.Continent{})
		db.Migrator().CreateTable(&model.Country{})
		db.Migrator().CreateTable(&model.City{})
		db.Migrator().CreateTable(&model.DeviceInfo{})
		db.Migrator().CreateTable(&model.DocumentFile{})
	}

	return db, nil
}
