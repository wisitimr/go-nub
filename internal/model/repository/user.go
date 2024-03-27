package repository

import (
	"context"
	mUser "findigitalservice/internal/model/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mUser.User, error)
	FindById(ctx context.Context, id string) (mUser.UserCompany, error)
	FindUserProfile(ctx context.Context, id primitive.ObjectID) (mUser.UserProfile, error)
	FindUserCompany(ctx context.Context, id primitive.ObjectID) (mUser.UserCompany, error)
	FindByUsername(ctx context.Context, email string) (mUser.User, error)
	FindByEmail(ctx context.Context, email string) (mUser.User, error)
	Create(ctx context.Context, payload mUser.User) (mUser.User, error)
	Update(ctx context.Context, payload mUser.User) (mUser.User, error)
}
