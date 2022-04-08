package domain

type Data interface {
	IsDataValid() error
	DataDomain() interface{}
}

const separator = "#"

type DynamoDomain struct {
	ID              string `dynamo:"ID,hash" index:"Prim-ID-index,hash"`
	IDType          string `dynamo:"IDTYPE"`
	PartitionID     string `dynamo:"PID,hash" index:"Sec-PID-index,range"`
	PartitionIDType string `dynamo:"PIDTYPE"`
}

func (domain *DynamoDomain) WithIndexes(id string, idType string, partitionID string, partitionIDType string) *DynamoDomain {
	domain.ID = idType + separator + id + separator + partitionIDType + separator + partitionID
	domain.IDType = idType
	domain.PartitionID = partitionIDType + separator + partitionID
	domain.PartitionIDType = partitionIDType
	return domain
}
