package repository

import (
	"context"
	mAccountType "nub/internal/model/account_type"
)

type AccountTypeRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mAccountType.AccountType, error)
	FindById(ctx context.Context, id string) (mAccountType.AccountType, error)
	Create(ctx context.Context, payload mAccountType.AccountType) (mAccountType.AccountType, error)
	Update(ctx context.Context, payload mAccountType.AccountType) (mAccountType.AccountType, error)
	Delete(ctx context.Context, id string) error
}
