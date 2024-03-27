package repository

import (
	"context"
	mCompany "findigitalservice/internal/model/company"
)

type CompanyRepository interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mCompany.Company, error)
	FindById(ctx context.Context, id string) (mCompany.Company, error)
	Create(ctx context.Context, payload mCompany.Company) (mCompany.Company, error)
	Update(ctx context.Context, payload mCompany.Company) (mCompany.Company, error)
}
