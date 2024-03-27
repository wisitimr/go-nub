package service

import (
	"context"
	mRole "findigitalservice/internal/model/role"
)

type RoleService interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mRole.Role, error)
	FindById(ctx context.Context, id string) (mRole.Role, error)
	Create(ctx context.Context, payload mRole.Role) (mRole.Role, error)
	Update(ctx context.Context, id string, payload mRole.Role) (mRole.Role, error)
}
