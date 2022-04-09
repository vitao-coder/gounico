package dynamodb

import "gounico/pkg/dynamodb/domain"

type DynamoClient interface {
	PutBatch(batchArray []interface{}) error
	Put(dynamoData domain.Data) error
	GetByIDAndPID(id string, pid string) (*domain.DynamoDomain, error)
	GetByID(id string) (*domain.DynamoDomain, error)
	GetByPID(pid string) (*[]domain.DynamoDomain, error)
	Update(id string, dynamoData domain.Data) error
	DeleteByID(id string) error
	DeleteByPID(id string) error
}
