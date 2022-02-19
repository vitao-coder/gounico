package feiralivre

import (
	"context"
	"gounico/feiralivre/domain"
	"gounico/pkg/errors"
)

type FeiraLivre interface {
	BuscarFeiraPorBairro(ctx context.Context, bairro string) ([]*domain.Feira, *errors.ServiceError)
}
