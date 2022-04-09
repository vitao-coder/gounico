package dynamodb

import "gounico/pkg/dynamodb/domain"

type DynamoClient interface {
	Put(dynamoData domain.Data) error
	GetByIDAndPID(id string, pid string, result *interface{}) error
	GetByID(id string, outResult *interface{}) error
	GetByPID(pid string) (*[]domain.DynamoDomain, error)
	Update(id string, dynamoData domain.Data) error
	DeleteByID(id string) error
	DeleteByPID(id string) error
}
