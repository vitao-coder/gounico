package feiralivre

import (
	"context"
	"gounico/feiralivre/domain"
	"gounico/pkg/errors"
)

type FeiraLivre interface {
	SaveFeira(ctx context.Context, request *domain.FeiraRequest) *errors.ServiceError
	ExcluirFeira(ctx context.Context, feiraID string, distritoID string) *errors.ServiceError
	BuscarFeiraPorDistrito(ctx context.Context, distritoID string) ([]domain.Feira, *errors.ServiceError)
}
