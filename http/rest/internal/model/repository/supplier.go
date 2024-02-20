package repository

import (
	"context"
	mSupplier "saved/http/rest/internal/model/supplier"
)

type SupplierRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mSupplier.Supplier, error)
	FindById(ctx context.Context, id string) (mSupplier.Supplier, error)
	Create(ctx context.Context, payload mSupplier.Supplier) (mSupplier.Supplier, error)
	Update(ctx context.Context, payload mSupplier.Supplier) (mSupplier.Supplier, error)
	Delete(ctx context.Context, id string) error
}
