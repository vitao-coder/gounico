package client

import (
	"gounico/pkg/dynamodb/domain"

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

func (dbc *DynamoDBClient) Put(dynamoData domain.Data) error {

	if err := dynamoData.IsDataValid(); err != nil {
		return err
	}
	if errPut := dbc.table.Put(dynamoData.DataDomain()).Run(); errPut != nil {
		return errPut
	}
	return nil
}

func (dbc *DynamoDBClient) GetByIDAndPID(id string, pid string) (*domain.DynamoDomain, error) {
	domain := &domain.DynamoDomain{}
	if err := dbc.table.Get("ID", id).
		Range("PID", dynamo.Equal, pid).
		One(&domain); err != nil {
		return nil, err
	}
	return domain, nil
}

func (dbc *DynamoDBClient) GetByID(id string) (*domain.DynamoDomain, error) {
	domain := &domain.DynamoDomain{}

	if err := dbc.table.Get("PRID", id).
		One(&domain); err != nil {
		return nil, err
	}
	return domain, nil
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

func (dbc *DynamoDBClient) DeleteByID(id string) error {
	if errUpd := dbc.table.Delete("ID", id).Run(); errUpd != nil {
		return errUpd
	}
	return nil
}

func (dbc *DynamoDBClient) DeleteByPID(id string) error {
	if errUpd := dbc.table.Delete("PID", id).Run(); errUpd != nil {
		return errUpd
	}
	return nil
}
