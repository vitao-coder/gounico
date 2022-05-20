package client

import (
	"context"
	"fmt"
	"gounico/pkg/database/dynamodb/domain"
	"gounico/pkg/telemetry/openTelemetry"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type DynamoDBClient struct {
	db    *dynamo.DB
	table *dynamo.Table
}

func NewDynamoDBClient(endpoint string, region string, id string, secret string, token string, tableName string) *DynamoDBClient {
	sess := session.Must(session.NewSession())
	cred := credentials.NewStaticCredentials(id, secret, token)
	db := dynamo.New(sess, &aws.Config{
		Credentials: cred,
		Region:      &region,
		Endpoint:    &endpoint,
	})

	db.Table(tableName).DeleteTable().Run()

	tableDyn := db.CreateTable(tableName, &domain.DynamoDomain{}).OnDemand(true)
	errCreateTable := tableDyn.Run()
	if errCreateTable != nil {
		panic(errCreateTable)
	}

	table := db.Table(tableName)
	return &DynamoDBClient{db: db, table: &table}
}

func (dbc *DynamoDBClient) PutBatch(batchArray []interface{}) error {
	_, err := dbc.table.Batch().Write().Put(batchArray...).Run()
	if err != nil {
		return err
	}

	return nil
}

func (dbc *DynamoDBClient) Put(ctx context.Context, dynamoData domain.Data) error {
	ctx, traceSpan := openTelemetry.TraceContextSpan(ctx, "DynamoDBClient - Put")
	defer traceSpan.End()

	if err := dynamoData.IsDataValid(); err != nil {
		return err
	}
	if errPut := dbc.table.Put(dynamoData.DataDomain()).Run(); errPut != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", errPut.Error()))
		openTelemetry.AddSpanError(traceSpan, errPut)
		return errPut
	}
	openTelemetry.SuccessSpan(traceSpan, "Success")
	return nil
}

func (dbc *DynamoDBClient) PutAsync(ctx context.Context, dynamoData domain.Data) error {
	var err error
	go func() {
		err = dbc.Put(ctx, dynamoData)
	}()
	return err
}

func (dbc *DynamoDBClient) GetByIDAndPID(id string, pid string) (*domain.DynamoDomain, error) {
	domainR := &domain.DynamoDomain{}

	if err := dbc.table.Get("PID", pid).
		Range("ID", dynamo.Equal, id).
		One(&domainR); err != nil {
		return nil, err
	}
	return domainR, nil
}

func (dbc *DynamoDBClient) GetByPID(pid string) (*[]domain.DynamoDomain, error) {

	domains := &[]domain.DynamoDomain{}

	if err := dbc.table.Get("PID", pid).
		All(domains); err != nil {
		return nil, err
	}
	return domains, nil
}

func (dbc *DynamoDBClient) Update(id string, dynamoData domain.Data) error {
	if err := dynamoData.IsDataValid(); err != nil {
		return err
	}
	if errUpd := dbc.table.Update(id, dynamoData.DataDomain()).Run(); errUpd != nil {
		return errUpd
	}
	return nil
}

func (dbc *DynamoDBClient) DeleteByIDAndPID(id string, pid string) error {
	if err := dbc.table.Delete("PID", pid).Range("ID", id).Run(); err != nil {
		return err
	}
	return nil
}
