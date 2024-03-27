package repository

import (
	"context"
	mProduct "findigitalservice/internal/model/product"
)

type ProductRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mProduct.Product, error)
	FindById(ctx context.Context, id string) (mProduct.Product, error)
	Create(ctx context.Context, payload mProduct.Product) (mProduct.Product, error)
	Update(ctx context.Context, payload mProduct.Product) (mProduct.Product, error)
	Delete(ctx context.Context, id string) error
}
