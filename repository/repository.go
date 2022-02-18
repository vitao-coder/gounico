package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	DB() *gorm.DB
	AutoMigrateDatabase() error
	Insert(ctx context.Context, model interface{}, createObject interface{}) error
	Save(saveObject *interface{}) error
	BulkInsert(ctx context.Context, model interface{}, createObjects ...interface{}) error
	Find(objectToFill *interface{}, fieldName string, fieldValue interface{}) error
}
