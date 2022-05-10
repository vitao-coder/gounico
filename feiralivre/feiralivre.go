package feiralivre

import (
	"context"
	"gounico/feiralivre/domains"
	"gounico/pkg/errors"
)

type FeiraLivre interface {
	SaveFeira(ctx context.Context, request *domains.FeiraRequest) *errors.ServiceError
	ExcluirFeira(ctx context.Context, feiraID string, distritoID string) *errors.ServiceError
	BuscarFeiraPorDistrito(ctx context.Context, distritoID string) ([]domains.Feira, *errors.ServiceError)
}
