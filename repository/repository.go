package repository

import (
	"gorm.io/gorm"
)

type Repository interface {
	DB() *gorm.DB
	AutoMigrateDatabase() error
}
