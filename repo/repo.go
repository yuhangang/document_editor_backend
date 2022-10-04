package repo

import (
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{db}
}
