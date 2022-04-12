package domain

type Data interface {
	IsDataValid() error
	DataDomain() interface{}
}

const Separator = "#"

type DynamoDomain struct {
	PartitionID   string `dynamo:"PID,hash" index:"Seq-ID-index,range"`
	PrimaryID     string `dynamo:"ID,range"`
	ID            string `dynamo:"PRID" localIndex:"ID-Seq-index,range" index:"Seq-ID-index,hash"`
	PrimaryType   string
	PartitionType string
	Data          interface{}
}

func NewDomainIndexes(id string, idType string, partitionID string, partitionIDType string) *DynamoDomain {
	return &DynamoDomain{
		ID:            partitionIDType + Separator + partitionID + idType + Separator + id,
		PrimaryID:     idType + Separator + id,
		PrimaryType:   idType,
		PartitionID:   partitionIDType + Separator + partitionID,
		PartitionType: partitionIDType,
	}
}
