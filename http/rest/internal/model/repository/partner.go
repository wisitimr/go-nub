package repository

import (
	"context"
	mPartner "findigitalservice/http/rest/internal/model/partner"
)

type PartnerRepository interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mPartner.Partner, error)
	FindById(ctx context.Context, id string) (mPartner.Partner, error)
	Create(ctx context.Context, payload mPartner.Partner) (mPartner.Partner, error)
	Update(ctx context.Context, payload mPartner.Partner) (mPartner.Partner, error)
}
