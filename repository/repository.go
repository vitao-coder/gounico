package repository

type Repository interface {
	AutoMigrate(migrateObject interface{}) error
	Insert(createObject *interface{}) error
	Save(saveObject *interface{}) error
	BulkInsert(createObjects ...interface{}) error
	Find(objectToFill *interface{}, fieldName string, fieldValue interface{}) error
}
