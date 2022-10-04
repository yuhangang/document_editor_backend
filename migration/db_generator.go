package migration

import (
	"echoapp/container"
	"echoapp/model"
)

// CreateDatabase creates the tables used in this application.
func CreateDatabase(container container.Container) {
	if container.GetConfig().Database.Migration {
		db := container.GetRepository()

		_ = db.DropTableIfExists(&model.Continent{})
		_ = db.DropTableIfExists(&model.Country{})
		_ = db.DropTableIfExists(&model.City{})

		_ = db.AutoMigrate(&model.Continent{})
		_ = db.AutoMigrate(&model.Country{})
		_ = db.AutoMigrate(&model.City{})
	}
}
