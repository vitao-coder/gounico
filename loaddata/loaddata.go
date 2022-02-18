package loaddata

import "gounico/pkg/errors"

type LoadData interface {
	ProcessCSVToDatabase(csvByteArray []byte) *errors.ServiceError
}
