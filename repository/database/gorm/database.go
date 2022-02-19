package gorm

import (
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

func (db *Database) DB() *gorm.DB {
	return db.GormDB
}
