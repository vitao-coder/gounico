package feiralivre

import (
	"context"
	"gounico/infrastructure/errors"
)

type ProcessCSV interface {
	ProcessCSVToDatabase(ctx context.Context, csvByteArray []byte) *errors.ServiceError
}
