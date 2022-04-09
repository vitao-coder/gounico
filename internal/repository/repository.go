package repository

import (
	client "gounico/pkg/dynamodb"
	"gounico/pkg/dynamodb/domain"
	"gounico/pkg/errors"
)

type Repository interface {
	Save(dynamoData domain.Data) *errors.ServiceError
	GetByPrimaryID(id string, entitytype string) (interface{}, *errors.ServiceError)
	GetBySecondaryID(id string, entitytype string) (*[]domain.DynamoDomain, *errors.ServiceError)
	Update(id string, dynamoData domain.Data) *errors.ServiceError
	Delete(id string) *errors.ServiceError
}

type repository struct {
	dataClient client.DynamoClient
}

func NewRepository(dataClient client.DynamoClient) *repository {
	return &repository{dataClient: dataClient}
}

func (repo *repository) Save(dynamoData domain.Data) *errors.ServiceError {
	err := repo.dataClient.Put(dynamoData)
	if err != nil {
		return errors.InternalServerError("error put data in dynamoDB", err)
	}
	return nil
}

func (repo *repository) GetByPrimaryID(id string, entitytype string) (interface{}, *errors.ServiceError) {
	return nil, nil
}

func (repo *repository) GetBySecondaryID(id string, entitytype string) (*[]domain.DynamoDomain, *errors.ServiceError) {

	feiraReturn, err := repo.dataClient.GetByPID(entitytype + domain.Separator + id)
	if err != nil {
		return nil, errors.InternalServerError("Error get by Secondary id", err)
	}

	return feiraReturn, nil
}

func (repo *repository) Update(id string, dynamoData domain.Data) *errors.ServiceError {
	return nil
}

func (repo *repository) Delete(id string) *errors.ServiceError {
	return nil
}
