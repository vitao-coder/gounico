package loaddata

import (
	"context"
	"gounico/pkg/errors"
)

type LoadData interface {
	ProcessCSVToDatabase(ctx context.Context, csvByteArray []byte) *errors.ServiceError
}
