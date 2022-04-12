package domain

type Data interface {
	IsDataValid() error
	DataDomain() interface{}
}

const Separator = "#"

type DynamoDomain struct {
	PartitionKey  string `dynamo:"PID,hash" index:"Seq-UID-index,range"`
	SortKey       string `dynamo:"ID,range"`
	UID           string `dynamo:"UID" localIndex:"UID-Seq-index,range" index:"Seq-UID-index,hash"`
	SortType      string
	PartitionType string
	Data          interface{}
}

func NewDomainIndexes(id string, idType string, partitionID string, partitionIDType string) *DynamoDomain {
	return &DynamoDomain{
		UID:           partitionIDType + Separator + partitionID + idType + Separator + id,
		SortKey:       idType + Separator + id,
		SortType:      idType,
		PartitionKey:  partitionIDType + Separator + partitionID,
		PartitionType: partitionIDType,
	}
}
