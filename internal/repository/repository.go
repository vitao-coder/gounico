package repository

import (
	client "gounico/pkg/dynamodb"
	"gounico/pkg/dynamodb/domain"
	"gounico/pkg/errors"
)

type Repository interface {
	BatchSave(dynamoData []interface{}) *errors.ServiceError
	Save(dynamoData domain.Data) *errors.ServiceError
	GetByID(id string, entitytype string) (interface{}, *errors.ServiceError)
	GetBySecondaryID(id string, entitytype string) (*[]domain.DynamoDomain, *errors.ServiceError)
	Delete(id string, entitytype string) *errors.ServiceError
}

type repository struct {
	dataClient client.DynamoClient
}

func NewRepository(dataClient client.DynamoClient) *repository {
	return &repository{dataClient: dataClient}
}

func (repo *repository) BatchSave(dynamoData []interface{}) *errors.ServiceError {
	err := repo.dataClient.PutBatch(dynamoData)
	if err != nil {
		return errors.InternalServerError("error put batch data in dynamoDB", err)
	}
	return nil
}

func (repo *repository) Save(dynamoData domain.Data) *errors.ServiceError {
	err := repo.dataClient.Put(dynamoData)
	if err != nil {
		return errors.InternalServerError("error put data in dynamoDB", err)
	}
	return nil
}

func (repo *repository) GetByID(id string, entitytype string) (interface{}, *errors.ServiceError) {
	feiraReturn, err := repo.dataClient.GetByID(entitytype + domain.Separator + id)
	if err != nil {
		return nil, errors.InternalServerError("Error get by id", err)
	}

	if feiraReturn == nil {
		return nil, errors.NotFoundError()
	}

	return feiraReturn, nil
}

func (repo *repository) GetBySecondaryID(id string, entitytype string) (*[]domain.DynamoDomain, *errors.ServiceError) {

	feiraReturn, err := repo.dataClient.GetByPID(entitytype + domain.Separator + id)
	if err != nil {
		return nil, errors.InternalServerError("Error get by Secondary id", err)
	}

	if feiraReturn == nil {
		return nil, errors.NotFoundError()
	}

	return feiraReturn, nil
}

func (repo *repository) Delete(id string, entitytype string) *errors.ServiceError {

	entity, err := repo.GetByID(id, entitytype)
	if err != nil {
		return err
	}

	if entity != nil {
		if err := repo.Delete(id, entitytype); err != nil {
			return err
		}
	} else {
		return errors.NotFoundError()
	}

	return nil
}
