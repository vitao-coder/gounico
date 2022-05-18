package dynamodb

import (
	"context"
	"gounico/pkg/database/dynamodb/domain"
)

type DynamoClient interface {
	PutBatch(batchArray []interface{}) error
	Put(ctx context.Context, dynamoData domain.Data) error
	GetByIDAndPID(id string, pid string) (*domain.DynamoDomain, error)
	GetByPID(pid string) (*[]domain.DynamoDomain, error)
	Update(id string, dynamoData domain.Data) error
	DeleteByIDAndPID(id string, pid string) error
}
