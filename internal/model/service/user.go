package service

import (
	"context"
	mRes "nub/internal/model/response"
	mUser "nub/internal/model/user"
)

type UserService interface {
	Count(ctx context.Context) (mRes.CountDto, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mUser.User, error)
	FindById(ctx context.Context, id string) (mUser.UserCompany, error)
	FindUserProfile(ctx context.Context) (mUser.UserProfile, error)
	FindUserCompany(ctx context.Context) (mUser.UserCompany, error)
	Create(ctx context.Context, payload mUser.User) (mUser.User, error)
	Update(ctx context.Context, id string, payload mUser.User) (mUser.UpdatedUserProfile, error)
	Login(ctx context.Context, payload mUser.Login) (mUser.UserProfile, error)
}
