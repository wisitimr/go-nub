package service

import (
	"context"
	mMaterial "findigitalservice/http/rest/internal/model/material"
	mRes "findigitalservice/http/rest/internal/model/response"
)

type MaterialService interface {
	Count(ctx context.Context) (mRes.CountDto, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mMaterial.Material, error)
	FindById(ctx context.Context, id string) (mMaterial.Material, error)
	Create(ctx context.Context, payload mMaterial.Material) (mMaterial.Material, error)
	Update(ctx context.Context, id string, payload mMaterial.Material) (mMaterial.Material, error)
	Delete(ctx context.Context, id string) error
}
