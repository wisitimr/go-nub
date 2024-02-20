package service

import (
	"context"
	mProduct "saved/http/rest/internal/model/product"
	mRes "saved/http/rest/internal/model/response"
)

type ProductService interface {
	Count(ctx context.Context) (mRes.CountDto, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mProduct.Product, error)
	FindById(ctx context.Context, id string) (mProduct.Product, error)
	Create(ctx context.Context, payload mProduct.Product) (mProduct.Product, error)
	Update(ctx context.Context, id string, payload mProduct.Product) (mProduct.Product, error)
	Delete(ctx context.Context, id string) error
}
