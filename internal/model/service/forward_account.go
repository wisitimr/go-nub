package service

import (
	"context"
	mForwardAccount "findigitalservice/internal/model/forward_account"
	mRes "findigitalservice/internal/model/response"
)

type ForwardAccountService interface {
	Count(ctx context.Context) (mRes.CountDto, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mForwardAccount.ForwardAccount, error)
	FindById(ctx context.Context, id string) (mForwardAccount.ForwardAccount, error)
	Create(ctx context.Context, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error)
	Update(ctx context.Context, id string, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error)
	Delete(ctx context.Context, id string) error
}
