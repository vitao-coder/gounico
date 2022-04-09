package feiralivre

import (
	"context"
	"gounico/pkg/errors"
)

type ProcessCSV interface {
	ProcessCSVToDatabase(ctx context.Context, csvByteArray []byte) *errors.ServiceError
}
