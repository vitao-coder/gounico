package repository

import (
	"context"
	"fmt"
	client "gounico/pkg/database/dynamodb"
	"gounico/pkg/database/dynamodb/domain"
	"gounico/pkg/errors"
	"gounico/pkg/telemetry/openTelemetry"
)

type Repository interface {
	BatchSave(dynamoData []interface{}) *errors.ServiceError
	Save(ctx context.Context, dynamoData domain.Data) *errors.ServiceError
	GetByIDAndBySecondaryID(id string, entitytype string, secondaryId string, secondaryEntitytype string) (interface{}, *errors.ServiceError)
	GetBySecondaryID(id string, entitytype string) (*[]domain.DynamoDomain, *errors.ServiceError)
	Delete(id string, entitytype string, secondaryId string, secondaryEntitytype string) *errors.ServiceError
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

func (repo *repository) Save(ctx context.Context, dynamoData domain.Data) *errors.ServiceError {
	ctx, traceSpan := openTelemetry.NewSpan(ctx, "Repository.Save")
	defer traceSpan.End()
	err := repo.dataClient.Put(ctx, dynamoData)
	if err != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
		openTelemetry.AddSpanError(traceSpan, err)
		return errors.InternalServerError("error put data in dynamoDB", err)
	}
	openTelemetry.SuccessSpan(traceSpan, "Success")
	return nil
}

func (repo *repository) GetByIDAndBySecondaryID(id string, entitytype string, secondaryId string, secondaryEntitytype string) (interface{}, *errors.ServiceError) {
	feiraReturn, err := repo.dataClient.GetByIDAndPID(entitytype+domain.Separator+id, secondaryEntitytype+domain.Separator+secondaryId)
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

func (repo *repository) Delete(id string, entitytype string, secondaryId string, secondaryEntitytype string) *errors.ServiceError {

	entity, err := repo.GetByIDAndBySecondaryID(id, entitytype, secondaryId, secondaryEntitytype)
	if err != nil {
		return err
	}

	if entity != nil {
		if err := repo.dataClient.DeleteByIDAndPID(entitytype+domain.Separator+id, secondaryEntitytype+domain.Separator+secondaryId); err != nil {
			return errors.InternalServerError("Error delete by id and PID", err)
		}
	} else {
		return errors.NotFoundError()
	}

	return nil
}
