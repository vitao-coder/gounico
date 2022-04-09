package domain

import "time"

type Data interface {
	IsDataValid() error
	DataDomain() interface{}
}

const Separator = "#"

type DynamoDomain struct {
	PartitionID   string    `dynamo:"PID,hash" index:"Seq-ID-index,range"`
	Time          time.Time `dynamo:",range"`
	PrimaryID     string    `localIndex:"ID-Seq-index,range" index:"Seq-ID-index,hash"`
	ID            string    `dynamo:"ID" index:"UUID-index,hash"`
	PrimaryType   string
	PartitionType string
	Data          interface{}
}

func NewDomainIndexes(id string, idType string, partitionID string, partitionIDType string) *DynamoDomain {
	return &DynamoDomain{
		ID:            idType + Separator + id + Separator + partitionIDType + Separator + partitionID,
		PrimaryID:     idType + Separator + id,
		PrimaryType:   idType,
		PartitionID:   partitionIDType + Separator + partitionID,
		PartitionType: partitionIDType,
		Time:          time.Now(),
	}
}
