package repository

import (
	"context"
	mForwardAccount "findigitalservice/internal/model/forward_account"
)

type ForwardAccountRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mForwardAccount.ForwardAccount, error)
	FindById(ctx context.Context, id string) (mForwardAccount.ForwardAccount, error)
	FindOne(ctx context.Context, query map[string][]string) (mForwardAccount.ForwardAccount, error)
	Create(ctx context.Context, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error)
	Update(ctx context.Context, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error)
	Delete(ctx context.Context, id string) error
}
