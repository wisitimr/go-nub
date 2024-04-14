package service

import (
	"context"
	mCompany "nub/internal/model/company"
)

type CompanyService interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mCompany.Company, error)
	FindById(ctx context.Context, id string) (mCompany.Company, error)
	Create(ctx context.Context, payload mCompany.Company) (mCompany.Company, error)
	Update(ctx context.Context, id string, payload mCompany.Company) (mCompany.Company, error)
}
