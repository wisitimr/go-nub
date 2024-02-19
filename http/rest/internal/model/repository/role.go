package repository

import (
	"context"
	mRole "saved/http/rest/internal/model/role"
)

type RoleRepository interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mRole.Role, error)
	FindById(ctx context.Context, id string) (mRole.Role, error)
	Create(ctx context.Context, payload mRole.Role) (mRole.Role, error)
	Update(ctx context.Context, payload mRole.Role) (mRole.Role, error)
}
