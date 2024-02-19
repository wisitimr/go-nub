package service

import (
	"context"
	mRes "saved/http/rest/internal/model/response"
	mSupplier "saved/http/rest/internal/model/supplier"
)

type SupplierService interface {
	Count(ctx context.Context) (mRes.CountDto, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mSupplier.Supplier, error)
	FindById(ctx context.Context, id string) (mSupplier.Supplier, error)
	Create(ctx context.Context, payload mSupplier.Supplier) (mSupplier.Supplier, error)
	Update(ctx context.Context, id string, payload mSupplier.Supplier) (mSupplier.Supplier, error)
}
