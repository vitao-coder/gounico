package feiralivre

import (
	"context"
	"gounico/feiralivre/domain"
	"gounico/pkg/errors"
)

type FeiraLivre interface {
	BuscarFeiraPorBairro(ctx context.Context, bairro string) ([]*domain.Feira, *errors.ServiceError)
	BuscarFeiraPorDistrito(ctx context.Context, distrito string) ([]*domain.Feira, *errors.ServiceError)
	ExcluirFeira(ctx context.Context, feiraID uint) *errors.ServiceError
	NovaFeira(ctx context.Context, request *domain.FeiraRequest) *errors.ServiceError
}
