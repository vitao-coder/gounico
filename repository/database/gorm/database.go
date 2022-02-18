package gorm

import (
	"context"
	"gounico/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	dbName string
	GormDB *gorm.DB
}

func NewDatabase(cfg config.Database) (*Database, error) {
	db, err := connectToDatabase(cfg)
	if err != nil {
		return nil, err
	}
	return &Database{
		dbName: cfg.Name,
		GormDB: db,
	}, nil
}

func connectToDatabase(cfg config.Database) (*gorm.DB, error) {
	dsn := cfg.Username + ":" + cfg.Password + "@tcp" + "(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Name + "?" + "parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) autoMigrate(migrateObject interface{}) error {

	if err := db.GormDB.AutoMigrate(migrateObject); err != nil {
		return err
	}
	return nil
}

func (db *Database) Insert(ctx context.Context, model interface{}, createObject interface{}) error {
	if err := db.GormDB.WithContext(ctx).Model(model).Create(&createObject); err != nil {
		return err.Error
	}
	return nil
}

func (db *Database) Save(saveObject *interface{}) error {
	if err := db.GormDB.Save(&saveObject); err != nil {
		return err.Error
	}
	return nil
}

func (db *Database) BulkInsert(ctx context.Context, model interface{}, createObjects ...interface{}) error {
	for _, object := range createObjects {
		db.Insert(ctx, model, object)
	}
	return nil
}

func (db *Database) Find(objectToFill *interface{}, fieldName string, fieldValue interface{}) error {
	if err := db.GormDB.Find(objectToFill, fieldName, fieldValue); err != nil {
		return err.Error
	}
	return nil
}

func (db *Database) DB() *gorm.DB {
	return db.GormDB
}
