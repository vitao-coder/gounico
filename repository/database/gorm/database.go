package gorm

import (
	"gorm.io/gorm/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	dbName string
	gormDB *gorm.DB
}

func NewDatabase(dbName string) (*Database, error) {
	db, err := connectToDatabase(dbName)
	if err != nil {
		return nil, err
	}
	return &Database{
		dbName: dbName,
		gormDB: db,
	}, nil
}

func connectToDatabase(dbname string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) AutoMigrate(migrateObject interface{}) error {

	if err := db.gormDB.AutoMigrate(migrateObject); err != nil {
		return err
	}
	return nil
}

func (db *Database) Insert(createObject *interface{}) error {
	if err := db.gormDB.Create(&createObject); err != nil {
		return err.Error
	}
	return nil
}

func (db *Database) Save(saveObject *interface{}) error {
	if err := db.gormDB.Save(&saveObject); err != nil {
		return err.Error
	}
	return nil
}

func (db *Database) BulkInsert(createObjects ...interface{}) error {
	if err := db.gormDB.Create(createObjects); err != nil {
		return err.Error
	}
	return nil
}

func (db *Database) Find(objectToFill *interface{}, fieldName string, fieldValue interface{}) error {
	if err := db.gormDB.Find(objectToFill, fieldName, fieldValue); err != nil {
		return err.Error
	}
	return nil
}
