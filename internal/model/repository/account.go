package repository

import (
	"context"
	mAccount "findigitalservice/internal/model/account"
)

type AccountRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mAccount.Account, error)
	FindById(ctx context.Context, id string) (mAccount.Account, error)
	Create(ctx context.Context, payload mAccount.Account) (mAccount.Account, error)
	Update(ctx context.Context, payload mAccount.Account) (mAccount.Account, error)
	Delete(ctx context.Context, id string) error
}
