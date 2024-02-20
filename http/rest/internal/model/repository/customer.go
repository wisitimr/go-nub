package repository

import (
	"context"
	mCustomer "saved/http/rest/internal/model/customer"
)

type CustomerRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mCustomer.Customer, error)
	FindById(ctx context.Context, id string) (mCustomer.Customer, error)
	Create(ctx context.Context, payload mCustomer.Customer) (mCustomer.Customer, error)
	Update(ctx context.Context, payload mCustomer.Customer) (mCustomer.Customer, error)
	Delete(ctx context.Context, id string) error
}
