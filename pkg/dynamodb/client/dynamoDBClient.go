package client

import (
	"gounico/pkg/dynamodb/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
	table := db.Table(tableName)
	return &DynamoDBClient{db: db, table: &table}
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

func (dbc *DynamoDBClient) GetByIDAndPID(id string, pid string, result *interface{}) error {

	if err := dbc.table.Get("ID", id).
		Range("PID", dynamo.Equal, pid).
		One(&result); err != nil {
		return err
	}
	return nil
}

func (dbc *DynamoDBClient) GetByID(id string, outResult *interface{}) error {

	if err := dbc.table.Get("ID", id).
		One(&outResult); err != nil {
		return err
	}
	return nil
}

func (dbc *DynamoDBClient) GetByPID(pid string, outResult []*interface{}) error {

	if err := dbc.table.Get("PID", pid).
		All(&outResult); err != nil {
		return err
	}
	return nil
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
