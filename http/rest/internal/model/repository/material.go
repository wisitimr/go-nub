package repository

import (
	"context"
	mMaterial "saved/http/rest/internal/model/material"
)

type MaterialRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mMaterial.Material, error)
	FindById(ctx context.Context, id string) (mMaterial.Material, error)
	Create(ctx context.Context, payload mMaterial.Material) (mMaterial.Material, error)
	Update(ctx context.Context, payload mMaterial.Material) (mMaterial.Material, error)
	Delete(ctx context.Context, id string) error
}
